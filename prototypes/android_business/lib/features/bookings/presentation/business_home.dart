import 'package:flutter/material.dart';

import '../../../app/theme.dart';
import '../../../local/demo_repository.dart';
import '../domain/booking.dart';
import 'booking_form.dart';
import 'booking_details.dart';

class BusinessHome extends StatefulWidget {
  const BusinessHome({
    super.key,
    required this.repository,
    required this.palette,
  });
  final DemoRepository repository;
  final UbiPalette palette;
  @override
  State<BusinessHome> createState() => _BusinessHomeState();
}

class _BusinessHomeState extends State<BusinessHome> {
  late DemoScope scope;
  late DateTime day;
  List<LocalBooking> bookings = [];
  int selectedTab = 0;
  int loadGeneration = 0;
  bool loading = true;
  String? error;

  @override
  void initState() {
    super.initState();
    final first = widget.repository.fixtures.operations.first;
    scope = DemoScope(first.organizationId, first.id);
    day = widget.repository.fixtures.referenceDate;
    _reload();
  }

  Future<void> _reload() async {
    final generation = ++loadGeneration;
    try {
      final result = await widget.repository.list(scope);
      if (!mounted || generation != loadGeneration) return;
      setState(() {
        bookings = result;
        loading = false;
        error = null;
      });
    } catch (_) {
      if (!mounted || generation != loadGeneration) return;
      setState(() {
        loading = false;
        error = 'The saved demo schedule could not load. Retry without adding it again.';
      });
    }
  }

  Future<void> _add() async {
    final result = await Navigator.of(context).push<LocalBooking>(
      MaterialPageRoute(
        builder: (_) => BookingForm(
          repository: widget.repository,
          palette: widget.palette,
          scope: scope,
          initialDay: day,
        ),
      ),
    );
    if (!mounted) return;
    if (result != null) {
      setState(() {
        day = widget.repository.fixtures.utcToLocal(result.startUtc);
        selectedTab = 0;
      });
    }
    await _reload();
  }

  Future<void> _open(LocalBooking booking) async {
    await Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => BookingDetails(
          repository: widget.repository,
          palette: widget.palette,
          booking: booking,
        ),
      ),
    );
    if (mounted) await _reload();
  }

  Future<void> _chooseDay() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: day,
      firstDate: DateTime(2029),
      lastDate: DateTime(2031, 12, 31),
      helpText: 'Demo date · Gregorian calendar',
    );
    if (picked != null && mounted) setState(() => day = picked);
  }

  @override
  Widget build(BuildContext context) {
    final p = widget.palette;
    final fixtures = widget.repository.fixtures;
    final org = fixtures.organizations.firstWhere(
      (o) => o.id == scope.organizationId,
    );
    final operation = fixtures.operations.firstWhere(
      (o) =>
          o.organizationId == scope.organizationId && o.id == scope.operationId,
    );
    final visible = bookings.where((b) {
      final local = fixtures.utcToLocal(b.startUtc);
      return local.year == day.year &&
          local.month == day.month &&
          local.day == day.day;
    }).toList();
    final pending = bookings
        .where((b) => b.status == LocalStatus.pending)
        .length;
    final conflicts = bookings
        .where((b) => b.status == LocalStatus.simulatedConflict)
        .length;
    return Scaffold(
      appBar: AppBar(
        title: const Text('Ubi Business'),
        actions: [
          IconButton(
            tooltip: 'About this synthetic demo',
            icon: const Icon(Icons.info_outline),
            onPressed: () => showDialog(
              context: context,
              builder: (_) => AlertDialog(
                title: const Text('A safe place to test'),
                content: const Text(
                  'This prototype uses invented businesses, labels and prices. Nothing is sent to a server. No messages, calls or payments work here. English copy and device performance still need validation.',
                ),
                actions: [
                  TextButton(
                    onPressed: () => Navigator.pop(context),
                    child: const Text('Understood'),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
      body: SafeArea(
        child: selectedTab == 2
            ? _lab(org.name, operation.name)
            : RefreshIndicator(
                onRefresh: _reload,
                child: ListView(
                  padding: const EdgeInsets.fromLTRB(20, 20, 20, 24),
                  children: [
                    DemoNotice(palette: p),
                    const SizedBox(height: 24),
                    Text(
                      org.name,
                      style: Theme.of(context).textTheme.headlineMedium,
                    ),
                    Text(
                      operation.name,
                      style: Theme.of(context).textTheme.bodyLarge,
                    ),
                    const SizedBox(height: 20),
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: p.infoSoft,
                        borderRadius: BorderRadius.circular(18),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Icon(Icons.cloud_off_outlined, color: p.info),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  'Local demo · no server',
                                  style: TextStyle(
                                    color: p.info,
                                    fontWeight: FontWeight.w700,
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 6),
                          Text(
                            'Saved items stay on this device. They are not confirmed bookings.',
                            style: TextStyle(color: p.info),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 20),
                    Wrap(
                      spacing: 12,
                      runSpacing: 12,
                      children: [
                        _stat(
                          'Saved locally',
                          pending.toString(),
                          Icons.save_outlined,
                          p.primarySoft,
                          p.primary,
                        ),
                        _stat(
                          'Needs review',
                          conflicts.toString(),
                          Icons.priority_high,
                          p.warningSoft,
                          p.warning,
                        ),
                      ],
                    ),
                    const SizedBox(height: 26),
                    Text(
                      selectedTab == 1
                          ? 'Choose a demo day'
                          : 'Your demo schedule',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 10),
                    OutlinedButton.icon(
                      onPressed: _chooseDay,
                      icon: const Icon(Icons.calendar_month_outlined),
                      label: Text('${demoDate(day)} · Gregorian'),
                    ),
                    const SizedBox(height: 8),
                    const Text('Addis Ababa · international 24-hour time'),
                    const SizedBox(height: 16),
                    if (loading) const LinearProgressIndicator(),
                    if (error != null) ...[
                      Text(error!, style: TextStyle(color: p.danger)),
                      OutlinedButton(
                        onPressed: _reload,
                        child: const Text('Retry loading'),
                      ),
                    ],
                    if (!loading && error == null && visible.isEmpty)
                      Container(
                        padding: const EdgeInsets.all(24),
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(20),
                          border: Border.all(color: p.border),
                        ),
                        child: Column(
                          children: [
                            Icon(
                              Icons.photo_camera_outlined,
                              size: 40,
                              color: p.primary,
                            ),
                            const SizedBox(height: 14),
                            Text(
                              'Start with one booking',
                              style: Theme.of(context).textTheme.titleMedium,
                            ),
                            const SizedBox(height: 8),
                            const Text(
                              'Try a call or walk-in using a made-up customer. Your first saved item will appear here.',
                              textAlign: TextAlign.center,
                            ),
                          ],
                        ),
                      ),
                    for (final booking in visible)
                      Padding(
                        padding: const EdgeInsets.only(bottom: 12),
                        child: BookingCard(
                          booking: booking,
                          fixtures: fixtures,
                          palette: p,
                          onTap: () => _open(booking),
                        ),
                      ),
                  ],
                ),
              ),
      ),
      // Reserve layout space for the primary action. A floating button can
      // cover a booking/status partway through scrolling, even with end padding.
      bottomNavigationBar: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (selectedTab != 2)
            DecoratedBox(
              decoration: BoxDecoration(
                color: Colors.white,
                border: Border(top: BorderSide(color: p.border)),
              ),
              child: SafeArea(
                top: false,
                bottom: false,
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(20, 12, 20, 12),
                  child: SizedBox(
                    width: double.infinity,
                    child: FilledButton.icon(
                      onPressed: _add,
                      icon: const Icon(Icons.add),
                      label: const Text(
                        'Add demo booking',
                        textAlign: TextAlign.center,
                      ),
                    ),
                  ),
                ),
              ),
            ),
          NavigationBar(
            selectedIndex: selectedTab,
            onDestinationSelected: (value) =>
                setState(() => selectedTab = value),
            destinations: const [
              NavigationDestination(
                icon: Icon(Icons.today_outlined),
                selectedIcon: Icon(Icons.today),
                label: 'Today',
              ),
              NavigationDestination(
                icon: Icon(Icons.calendar_month_outlined),
                selectedIcon: Icon(Icons.calendar_month),
                label: 'Calendar',
              ),
              NavigationDestination(
                icon: Icon(Icons.science_outlined),
                selectedIcon: Icon(Icons.science),
                label: 'Demo lab',
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _stat(
    String label,
    String value,
    IconData icon,
    Color background,
    Color foreground,
  ) => Container(
    width: 150,
    padding: const EdgeInsets.all(16),
    decoration: BoxDecoration(
      color: background,
      borderRadius: BorderRadius.circular(18),
    ),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(icon, color: foreground),
        const SizedBox(height: 10),
        Text(
          value,
          style: TextStyle(
            fontSize: 28,
            fontWeight: FontWeight.w700,
            color: foreground,
          ),
        ),
        Text(label, style: TextStyle(color: foreground)),
      ],
    ),
  );

  Widget _lab(String organization, String operation) => ListView(
    padding: const EdgeInsets.all(20),
    children: [
      DemoNotice(palette: widget.palette),
      const SizedBox(height: 24),
      Text(
        'Test the boundaries',
        style: Theme.of(context).textTheme.headlineMedium,
      ),
      const SizedBox(height: 8),
      const Text(
        'Switch the invented business or operation. Saved items stay in their own scope. This is not sign-in or a permission system.',
      ),
      const SizedBox(height: 20),
      for (final option in widget.repository.fixtures.operations)
        Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: OutlinedButton(
            onPressed: () {
              setState(() {
                scope = DemoScope(option.organizationId, option.id);
                bookings = [];
                loading = true;
              });
              _reload();
            },
            child: Text(
              '${widget.repository.fixtures.organizations.firstWhere((o) => o.id == option.organizationId).name}\n${option.name}',
            ),
          ),
        ),
      const SizedBox(height: 16),
      Text('Selected: $organization · $operation'),
      const SizedBox(height: 24),
      const Text(
        'To test a conflict, open a saved item and choose “Simulate a conflict.” No server is contacted. To test persistence, close and reopen the demo.',
      ),
      const SizedBox(height: 24),
      const Text(
        'Not built in V0: real account access, server synchronization, quotation approval, messaging, calls, payments, Amharic UI and Ethiopian-calendar entry. These remain planned work, not working demo buttons.',
      ),
    ],
  );
}

class DemoNotice extends StatelessWidget {
  const DemoNotice({super.key, required this.palette});
  final UbiPalette palette;
  @override
  Widget build(BuildContext context) => Container(
    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
    decoration: BoxDecoration(
      color: palette.primarySoft,
      borderRadius: BorderRadius.circular(12),
    ),
    child: Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(Icons.science_outlined, color: palette.primary, size: 20),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            'SYNTHETIC DEMO · no real bookings or money',
            style: TextStyle(
              color: palette.primary,
              fontSize: 13,
              fontWeight: FontWeight.w700,
            ),
          ),
        ),
      ],
    ),
  );
}

class BookingCard extends StatelessWidget {
  const BookingCard({
    super.key,
    required this.booking,
    required this.fixtures,
    required this.palette,
    required this.onTap,
  });
  final LocalBooking booking;
  final DemoFixtures fixtures;
  final UbiPalette palette;
  final VoidCallback onTap;
  @override
  Widget build(BuildContext context) {
    final conflict = booking.status == LocalStatus.simulatedConflict;
    return Material(
      color: Colors.white,
      borderRadius: BorderRadius.circular(18),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(18),
        child: Padding(
          padding: const EdgeInsets.all(18),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '${demoTime(fixtures.utcToLocal(booking.startUtc))} – ${demoTime(fixtures.utcToLocal(booking.endUtc))}',
                style: TextStyle(
                  color: palette.primary,
                  fontWeight: FontWeight.w700,
                  fontSize: 16,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                booking.serviceName,
                style: Theme.of(context).textTheme.titleMedium,
              ),
              Text('${booking.customerLabel} · ${booking.source}'),
              const SizedBox(height: 12),
              Text(
                conflict
                    ? 'Simulated conflict · review needed'
                    : 'Saved on this phone · not confirmed',
                style: TextStyle(
                  color: conflict ? palette.warning : palette.info,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 8),
              const Text(
                'Open saved details',
                style: TextStyle(decoration: TextDecoration.underline),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
