import 'dart:convert';
import 'dart:math';

enum LocalStatus { pending, simulatedConflict, cancelled }

class DemoScope {
  const DemoScope(this.organizationId, this.operationId);
  final String organizationId;
  final String operationId;
  String get key => '$organizationId/$operationId';
}

class DemoOrganization {
  DemoOrganization(Map<String, dynamic> json)
    : id = json['id'] as String,
      name = json['name'] as String;
  final String id;
  final String name;
}

class DemoOperation {
  DemoOperation(Map<String, dynamic> json)
    : organizationId = json['organizationId'] as String,
      id = json['id'] as String,
      name = json['name'] as String;
  final String organizationId;
  final String id;
  final String name;
}

class DemoService {
  DemoService(Map<String, dynamic> json)
    : organizationId = json['organizationId'] as String,
      operationId = json['operationId'] as String,
      id = json['id'] as String,
      name = json['name'] as String,
      durationMinutes = json['durationMinutes'] as int,
      priceMinor = json['priceMinor'] as int,
      resourceId = json['resourceId'] as String;
  final String organizationId;
  final String operationId;
  final String id;
  final String name;
  final int durationMinutes;
  final int priceMinor;
  final String resourceId;
}

class DemoFixtures {
  DemoFixtures._(Map<String, dynamic> json)
    : referenceDate = DateTime.parse(json['referenceDate'] as String),
      offsetMinutes = json['utcOffsetMinutes'] as int,
      organizations = (json['organizations'] as List)
          .map((e) => DemoOrganization(e as Map<String, dynamic>))
          .toList(),
      operations = (json['operations'] as List)
          .map((e) => DemoOperation(e as Map<String, dynamic>))
          .toList(),
      services = (json['services'] as List)
          .map((e) => DemoService(e as Map<String, dynamic>))
          .toList();

  factory DemoFixtures.parse(String source) {
    final json = jsonDecode(source) as Map<String, dynamic>;
    if (json['syntheticOnly'] != true ||
        json['schemaVersion'] != 1 ||
        json['timezone'] != 'Africa/Addis_Ababa' ||
        json['utcOffsetMinutes'] != 180 ||
        json['currency'] != 'ETB' ||
        json['currencyExponent'] != 2 ||
        json['clockConvention'] != 'international-24h') {
      throw const FormatException('Unsupported or non-synthetic fixtures.');
    }
    final fixtures = DemoFixtures._(json);
    final scopes = fixtures.operations
        .map((o) => '${o.organizationId}/${o.id}')
        .toSet();
    if (scopes.length != fixtures.operations.length ||
        fixtures.services.isEmpty) {
      throw const FormatException('Invalid fixture scopes.');
    }
    final serviceKeys = <String>{};
    for (final operation in fixtures.operations) {
      if (!fixtures.organizations.any(
        (o) => o.id == operation.organizationId,
      )) {
        throw const FormatException('Unknown fixture organization.');
      }
    }
    for (final service in fixtures.services) {
      final scope = '${service.organizationId}/${service.operationId}';
      if (!scopes.contains(scope) ||
          service.durationMinutes <= 0 ||
          service.priceMinor < 0 ||
          service.resourceId.isEmpty ||
          !serviceKeys.add('$scope/${service.id}')) {
        throw const FormatException('Invalid fixture service.');
      }
    }
    return fixtures;
  }

  final DateTime referenceDate;
  final int offsetMinutes;
  final List<DemoOrganization> organizations;
  final List<DemoOperation> operations;
  final List<DemoService> services;
  List<DemoService> servicesFor(DemoScope scope) => services
      .where(
        (s) =>
            s.organizationId == scope.organizationId &&
            s.operationId == scope.operationId,
      )
      .toList();
  DemoService serviceFor(DemoScope scope, String id) => servicesFor(scope)
      .firstWhere(
        (s) => s.id == id,
        orElse: () =>
            throw ArgumentError('Service is outside this demo scope.'),
      );
  DateTime localToUtc(DateTime date, int hour, int minute) => DateTime.utc(
    date.year,
    date.month,
    date.day,
    hour,
    minute,
  ).subtract(Duration(minutes: offsetMinutes));
  DateTime utcToLocal(DateTime utc) =>
      utc.toUtc().add(Duration(minutes: offsetMinutes));
}

class BookingIntent {
  const BookingIntent({
    required this.scope,
    required this.serviceId,
    required this.customerLabel,
    required this.startUtc,
    required this.source,
  });
  final DemoScope scope;
  final String serviceId;
  final String customerLabel;
  final DateTime startUtc;
  final String source;
  Map<String, Object?> toJson() => {
    'organizationId': scope.organizationId,
    'operationId': scope.operationId,
    'serviceId': serviceId,
    'customerLabel': customerLabel.trim(),
    'startUtcMs': startUtc.toUtc().millisecondsSinceEpoch,
    'source': source,
  };
  factory BookingIntent.fromJson(Map<String, dynamic> json) => BookingIntent(
    scope: DemoScope(
      json['organizationId'] as String,
      json['operationId'] as String,
    ),
    serviceId: json['serviceId'] as String,
    customerLabel: json['customerLabel'] as String,
    startUtc: DateTime.fromMillisecondsSinceEpoch(
      json['startUtcMs'] as int,
      isUtc: true,
    ),
    source: json['source'] as String,
  );
}

class LocalBooking {
  LocalBooking.fromRow(Map<String, Object?> row)
    : id = row['id'] as String,
      scope = DemoScope(
        row['organization_id'] as String,
        row['operation_id'] as String,
      ),
      serviceId = row['service_id'] as String,
      serviceName = row['service_name'] as String,
      customerLabel = row['customer_label'] as String,
      source = row['source'] as String,
      startUtc = DateTime.fromMillisecondsSinceEpoch(
        row['start_ms'] as int,
        isUtc: true,
      ),
      endUtc = DateTime.fromMillisecondsSinceEpoch(
        row['end_ms'] as int,
        isUtc: true,
      ),
      priceMinor = row['price_minor'] as int,
      status = LocalStatus.values.byName(row['status'] as String),
      revision = row['revision'] as int;
  final String id;
  final DemoScope scope;
  final String serviceId;
  final String serviceName;
  final String customerLabel;
  final String source;
  final DateTime startUtc;
  final DateTime endUtc;
  final int priceMinor;
  final LocalStatus status;
  final int revision;
}

class LocalBookingConflict implements Exception {
  const LocalBookingConflict();
  @override
  String toString() =>
      'This phone already has a demo item at this time. Choose another time. Your details are kept.';
}

String newDemoId() {
  final random = Random.secure();
  return 'v0_${List.generate(16, (_) => random.nextInt(256).toRadixString(16).padLeft(2, '0')).join()}';
}

String demoMoney(int minor) =>
    'ETB ${minor ~/ 100}.${(minor % 100).toString().padLeft(2, '0')}';

String demoDate(DateTime date) {
  const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];
  return '${date.day} ${months[date.month - 1]} ${date.year}';
}

String demoTime(DateTime time) =>
    '${time.hour.toString().padLeft(2, '0')}:${time.minute.toString().padLeft(2, '0')}';
