import 'dart:async';

import 'package:flutter/material.dart';

import '../../../app/theme.dart';
import '../../../local/demo_repository.dart';
import '../domain/booking.dart';

class BookingForm extends StatefulWidget {
  const BookingForm({
    super.key,
    required this.repository,
    required this.palette,
    required this.scope,
    required this.initialDay,
    this.existing,
  });
  final DemoRepository repository;
  final UbiPalette palette;
  final DemoScope scope;
  final DateTime initialDay;
  final LocalBooking? existing;
  @override
  State<BookingForm> createState() => _BookingFormState();
}

class _BookingFormState extends State<BookingForm> {
  late String serviceId;
  String customer = 'Demo Guest 01';
  String source = 'Call';
  late DateTime day;
  int hour = 9;
  int minute = 0;
  bool loading = true;
  bool saving = false;
  bool reviewing = false;
  bool draftSaved = false;
  String? error;
  int draftGeneration = 0;
  Future<void> draftWrites = Future.value();
  final commandId = newDemoId();

  BookingIntent get intent => BookingIntent(
    scope: widget.scope,
    serviceId: serviceId,
    customerLabel: customer,
    source: source,
    startUtc: widget.repository.fixtures.localToUtc(day, hour, minute),
  );

  @override
  void initState() {
    super.initState();
    day = widget.initialDay;
    serviceId = widget.repository.fixtures.servicesFor(widget.scope).first.id;
    _load();
  }

  Future<void> _load() async {
    try {
      final existing = widget.existing;
      final draft = existing == null
          ? await widget.repository.readDraft(widget.scope)
          : BookingIntent(
              scope: existing.scope,
              serviceId: existing.serviceId,
              customerLabel: existing.customerLabel,
              source: existing.source,
              startUtc: existing.startUtc,
            );
      if (draft != null) {
        final local = widget.repository.fixtures.utcToLocal(draft.startUtc);
        serviceId = draft.serviceId;
        customer = draft.customerLabel;
        source = draft.source;
        day = local;
        hour = local.hour;
        minute = local.minute;
      }
      if (!mounted) return;
      setState(() {
        loading = false;
        draftSaved = draft != null && existing == null;
      });
    } catch (_) {
      if (!mounted) return;
      setState(() {
        loading = false;
        error =
            'The previous draft could not load. No booking has been submitted.';
      });
    }
  }

  void _changed(VoidCallback update) {
    setState(() {
      update();
      reviewing = false;
      draftSaved = false;
      error = null;
    });
    final snapshot = intent;
    final generation = ++draftGeneration;
    draftWrites = draftWrites
        .then((_) => widget.repository.saveDraft(snapshot))
        .then((_) {
          if (mounted && generation == draftGeneration) {
            setState(() => draftSaved = true);
          }
        })
        .catchError((Object _) {
          if (mounted && generation == draftGeneration) {
            setState(() {
              error = 'This draft was not saved. Check the time or retry. Your current choices remain on screen.';
            });
          }
        });
  }

  Future<void> _date() async {
    final result = await showDatePicker(
      context: context,
      initialDate: day,
      firstDate: DateTime(2029),
      lastDate: DateTime(2031, 12, 31),
      helpText: 'Demo date · Gregorian calendar',
    );
    if (result != null && mounted) _changed(() => day = result);
  }

  Future<void> _review() async {
    await draftWrites;
    if (!mounted) return;
    setState(() {
      reviewing = true;
      error = null;
    });
  }

  Future<void> _save() async {
    if (saving) return;
    setState(() {
      saving = true;
      error = null;
    });
    try {
      await draftWrites;
      final result = widget.existing == null
          ? await widget.repository.save(intent, commandId: commandId)
          : await widget.repository.revise(
              widget.existing!,
              intent,
              commandId: commandId,
            );
      if (mounted) Navigator.of(context).pop(result);
    } on LocalBookingConflict catch (e) {
      if (mounted) {
        setState(() {
          saving = false;
          reviewing = false;
          error = e.toString();
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() {
          saving = false;
          error = 'The save could not finish. Nothing is confirmed. Retry the same save; do not add it again.';
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final service = widget.repository.fixtures.serviceFor(
      widget.scope,
      serviceId,
    );
    final operation = widget.repository.fixtures.operations.firstWhere(
      (o) =>
          o.organizationId == widget.scope.organizationId &&
          o.id == widget.scope.operationId,
    );
    final org = widget.repository.fixtures.organizations.firstWhere(
      (o) => o.id == widget.scope.organizationId,
    );
    final end = widget.repository.fixtures
        .utcToLocal(intent.startUtc)
        .add(Duration(minutes: service.durationMinutes));
    return PopScope(
      canPop: !saving,
      child: Scaffold(
        appBar: AppBar(
          title: Text(
            reviewing
                ? 'Review before saving'
                : widget.existing == null
                ? 'Add demo booking'
                : 'Change demo time',
          ),
        ),
        body: SafeArea(
          child: loading
              ? const Center(child: CircularProgressIndicator())
              : ListView(
                  padding: const EdgeInsets.all(20),
                  children: [
                    Text(
                      '${org.name} · ${operation.name}',
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      'Invented labels and prices only. No real customer information.',
                    ),
                    const SizedBox(height: 24),
                    if (error != null) ...[
                      Semantics(
                        liveRegion: true,
                        child: Text(
                          error!,
                          style: TextStyle(color: widget.palette.danger),
                        ),
                      ),
                      const SizedBox(height: 20),
                    ],
                    if (!reviewing) ...[
                      DropdownButtonFormField<String>(
                        key: ValueKey('service-$serviceId'),
                        initialValue: serviceId,
                        isExpanded: true,
                        decoration: const InputDecoration(labelText: 'Service'),
                        items: widget.repository.fixtures
                            .servicesFor(widget.scope)
                            .map(
                              (s) => DropdownMenuItem(
                                value: s.id,
                                child: Text(s.name),
                              ),
                            )
                            .toList(),
                        onChanged: (value) {
                          if (value != null) _changed(() => serviceId = value);
                        },
                      ),
                      const SizedBox(height: 16),
                      Text(
                        '${service.durationMinutes} minutes · demo price ${demoMoney(service.priceMinor)}',
                      ),
                      const SizedBox(height: 20),
                      OutlinedButton.icon(
                        onPressed: _date,
                        icon: const Icon(Icons.calendar_month_outlined),
                        label: Text('${demoDate(day)} · Gregorian'),
                      ),
                      const SizedBox(height: 20),
                      Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Expanded(
                            child: DropdownButtonFormField<int>(
                              key: ValueKey('hour-$hour'),
                              initialValue: hour,
                              decoration: const InputDecoration(
                                labelText: 'Hour · 24h',
                              ),
                              items: List.generate(
                                9,
                                (i) => DropdownMenuItem(
                                  value: i + 9,
                                  child: Text('${i + 9}'.padLeft(2, '0')),
                                ),
                              ),
                              onChanged: (value) {
                                if (value != null) _changed(() => hour = value);
                              },
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: DropdownButtonFormField<int>(
                              key: ValueKey('minute-$minute'),
                              initialValue: minute,
                              decoration: const InputDecoration(
                                labelText: 'Minute',
                              ),
                              items: const [
                                DropdownMenuItem(value: 0, child: Text('00')),
                                DropdownMenuItem(value: 30, child: Text('30')),
                              ],
                              onChanged: (value) {
                                if (value != null) {
                                  _changed(() => minute = value);
                                }
                              },
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      Text(
                        'Ends at ${demoTime(end)} · Addis Ababa, international time',
                      ),
                      const SizedBox(height: 20),
                      DropdownButtonFormField<String>(
                        key: ValueKey('customer-$customer'),
                        initialValue: customer,
                        isExpanded: true,
                        decoration: const InputDecoration(
                          labelText: 'Made-up customer',
                        ),
                        items:
                            ['Demo Guest 01', 'Demo Guest 02', 'Demo Group 03']
                                .map(
                                  (s) => DropdownMenuItem(
                                    value: s,
                                    child: Text(s),
                                  ),
                                )
                                .toList(),
                        onChanged: (value) {
                          if (value != null) _changed(() => customer = value);
                        },
                      ),
                      const SizedBox(height: 20),
                      DropdownButtonFormField<String>(
                        key: ValueKey('source-$source'),
                        initialValue: source,
                        isExpanded: true,
                        decoration: const InputDecoration(
                          labelText: 'Request source',
                        ),
                        items:
                            ['Call', 'Walk-in', 'WhatsApp', 'Telegram', 'Other']
                                .map(
                                  (s) => DropdownMenuItem(
                                    value: s,
                                    child: Text(s),
                                  ),
                                )
                                .toList(),
                        onChanged: (value) {
                          if (value != null) _changed(() => source = value);
                        },
                      ),
                      const SizedBox(height: 16),
                      const Text(
                        'Selecting a source does not connect to that app.',
                      ),
                      if (draftSaved) ...[
                        const SizedBox(height: 16),
                        Text(
                          'Draft saved on this phone · not submitted',
                          style: TextStyle(color: widget.palette.info),
                        ),
                      ],
                    ] else ...[
                      _line(context, 'Service', service.name),
                      _line(
                        context,
                        'When',
                        '${demoDate(day)} · Gregorian\n${hour.toString().padLeft(2, '0')}:${minute.toString().padLeft(2, '0')} – ${demoTime(end)}',
                      ),
                      _line(
                        context,
                        'Clock and place',
                        'Addis Ababa · international 24-hour time',
                      ),
                      _line(context, 'Customer', customer),
                      _line(context, 'Source', source),
                      _line(
                        context,
                        'Invented price',
                        '${demoMoney(service.priceMinor)} · no money collected',
                      ),
                      Container(
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          color: widget.palette.infoSoft,
                          borderRadius: BorderRadius.circular(16),
                        ),
                        child: Text(
                          'Saving creates a local tentative item. No server checks availability, no customer is contacted and nothing is confirmed.',
                          style: TextStyle(color: widget.palette.info),
                        ),
                      ),
                      const SizedBox(height: 20),
                      OutlinedButton(
                        onPressed: saving
                            ? null
                            : () => setState(() => reviewing = false),
                        child: const Text('Edit these details'),
                      ),
                    ],
                  ],
                ),
        ),
        bottomNavigationBar: SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: FilledButton(
              onPressed: loading || saving
                  ? null
                  : reviewing
                  ? _save
                  : _review,
              child: Text(
                saving
                    ? 'Saving on this phone…'
                    : reviewing
                    ? 'Save on this phone'
                    : 'Review details',
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _line(BuildContext context, String label, String value) => Padding(
    padding: const EdgeInsets.only(bottom: 20),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: Theme.of(context).textTheme.bodyMedium),
        const SizedBox(height: 4),
        Text(value, style: Theme.of(context).textTheme.bodyLarge),
      ],
    ),
  );
}
