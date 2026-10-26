import 'package:flutter/material.dart';

import '../../../app/theme.dart';
import '../../../local/demo_repository.dart';
import '../domain/booking.dart';
import 'booking_form.dart';

class BookingDetails extends StatefulWidget {
  const BookingDetails({
    super.key,
    required this.repository,
    required this.palette,
    required this.booking,
  });
  final DemoRepository repository;
  final UbiPalette palette;
  final LocalBooking booking;
  @override
  State<BookingDetails> createState() => _BookingDetailsState();
}

class _BookingDetailsState extends State<BookingDetails> {
  late LocalBooking booking;
  List<Map<String, Object?>> history = [];
  bool busy = false;
  String? error;
  @override
  void initState() {
    super.initState();
    booking = widget.booking;
    _reload();
  }

  Future<void> _reload() async {
    try {
      final item = await widget.repository.get(booking.scope, booking.id);
      final events = await widget.repository.history(booking.scope, booking.id);
      if (mounted) {
        setState(() {
          booking = item;
          history = events;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(
          () =>
              error = 'Saved details could not reload. No new change was made.',
        );
      }
    }
  }

  Future<void> _status(LocalStatus status) async {
    if (busy) return;
    if (status == LocalStatus.cancelled) {
      final allowed = await showDialog<bool>(
        context: context,
        builder: (_) => AlertDialog(
          title: const Text('Cancel this local demo item?'),
          content: const Text(
            'Only this synthetic item is cancelled. No real booking, customer or money is affected. Its local history is retained.',
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: const Text('Keep item'),
            ),
            FilledButton(
              onPressed: () => Navigator.pop(context, true),
              child: const Text('Cancel local item'),
            ),
          ],
        ),
      );
      if (allowed != true || !mounted) return;
    }
    setState(() {
      busy = true;
      error = null;
    });
    try {
      await widget.repository.changeStatus(booking.scope, booking.id, status);
      await _reload();
    } catch (_) {
      if (mounted) {
        setState(
          () => error = 'The local change could not finish. Reopen details before trying again.',
        );
      }
    } finally {
      if (mounted) setState(() => busy = false);
    }
  }

  Future<void> _edit() async {
    await Navigator.of(context).push<LocalBooking>(
      MaterialPageRoute(
        builder: (_) => BookingForm(
          repository: widget.repository,
          palette: widget.palette,
          scope: booking.scope,
          initialDay: widget.repository.fixtures.utcToLocal(booking.startUtc),
          existing: booking,
        ),
      ),
    );
    if (mounted) await _reload();
  }

  @override
  Widget build(BuildContext context) {
    final local = widget.repository.fixtures.utcToLocal(booking.startUtc);
    final conflict = booking.status == LocalStatus.simulatedConflict;
    final cancelled = booking.status == LocalStatus.cancelled;
    final status = cancelled
        ? 'Cancelled locally'
        : conflict
        ? 'Simulated conflict · review needed'
        : 'Saved on this phone · not confirmed';
    final color = conflict ? widget.palette.warning : widget.palette.info;
    const eventLabels = {
      'local_intent_saved': 'Created on this phone',
      'simulated_conflict': 'Conflict simulation',
      'local_intent_revised': 'Details changed locally',
      'local_intent_cancelled': 'Cancelled locally',
    };
    return Scaffold(
      appBar: AppBar(title: const Text('Saved demo details')),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            Text(
              booking.serviceName,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 12),
            Semantics(
              liveRegion: true,
              child: Text(
                status,
                style: TextStyle(color: color, fontWeight: FontWeight.w700),
              ),
            ),
            const SizedBox(height: 20),
            Text(
              '${demoDate(local)} · Gregorian\n${demoTime(local)} – ${demoTime(widget.repository.fixtures.utcToLocal(booking.endUtc))}',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            const Text('Addis Ababa · international 24-hour time'),
            const SizedBox(height: 20),
            Text(
              '${booking.customerLabel}\nSource: ${booking.source}\nDemo price: ${demoMoney(booking.priceMinor)}',
            ),
            const SizedBox(height: 20),
            Text(
              conflict
                  ? 'This is a deliberately simulated conflict. No server was contacted. Your original details are kept; choose a different time to create a revised local intent.'
                  : 'This is not a confirmed booking. No server, customer message or payment is connected.',
            ),
            const SizedBox(height: 20),
            if (error != null) ...[
              Text(error!, style: TextStyle(color: widget.palette.danger)),
              const SizedBox(height: 20),
            ],
            if (!cancelled) ...[
              FilledButton(
                onPressed: busy ? null : _edit,
                child: Text(
                  conflict ? 'Choose another time' : 'Change local details',
                ),
              ),
              const SizedBox(height: 12),
              if (!conflict)
                OutlinedButton(
                  onPressed: busy
                      ? null
                      : () => _status(LocalStatus.simulatedConflict),
                  child: const Text('Simulate a conflict'),
                ),
              const SizedBox(height: 12),
              OutlinedButton(
                onPressed: busy ? null : () => _status(LocalStatus.cancelled),
                child: const Text('Cancel local demo item'),
              ),
            ],
            const SizedBox(height: 28),
            Text(
              'Local history',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 12),
            for (final event in history)
              Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: Text(
                  '• ${eventLabels[event['action']] ?? 'Local event'}',
                ),
              ),
            const SizedBox(height: 12),
            Text(
              'Demo reference: ${booking.id}',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          ],
        ),
      ),
    );
  }
}
