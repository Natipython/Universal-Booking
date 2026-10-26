# ADR-002 — Language, runtime and native-code strategy

Status: **Proposed — requires prototype, staffing and operational evidence**  
Decision owner: CTO  
Required approvers: product/mobile/backend/security/reliability owners and CEO/CFO for staffing/cost  
Date proposed: 2026-09-18 (Africa/Addis_Ababa)

## Problem

Ubi needs a fast Android-first merchant product, a lightweight customer web path, future iOS/desktop surfaces and a reliable provider-integrating backend. Using every capable language everywhere would multiply build systems, security surfaces, hiring needs, debugging paths and release risk.

The objective is therefore not “maximum number of languages.” It is the smallest maintainable set with explicit boundaries and room for C, C++, Rust, Java and Swift where they produce measured value.

## Proposed decision

| Area | Primary choice | Allowed supporting choice | Explicit boundary |
|---|---|---|---|
| Server application/API/workers/operational CLIs | Go | SQL; generated contracts | Modular monolith plus workers; no C/C++ business-domain fork |
| Ubi Business mobile | Flutter/Dart, Android first | Kotlin/Java Android plugin code; later Swift iOS plugin code | Native adapters expose platform capability, not separate booking rules |
| Customer web, business desktop web, operations console | React + TypeScript with server rendering/progressive enhancement | Standards-based HTML/CSS and generated API types | Critical customer content/action works before heavy client enhancement |
| Authoritative data | PostgreSQL SQL/migrations | Go data access and generated schema/contracts | Database invariants remain authoritative |
| New memory-sensitive native component | Rust preferred when team/tooling can support it | C ABI wrapper | Separate package/process/library with threat model and ownership |
| Existing/high-value native library or measured hotspot | C or C++ | Dart FFI, JNI or Swift interop through narrow wrapper | No raw unsafe/native interface spread through app code |
| Platform-specific Apple integration | Swift | Objective-C/C/C++ interop if vendor/library requires | Thin Flutter plugin/adapter |
| Platform-specific Android integration | Kotlin/Java | JNI to audited C/C++ only where required | Thin Flutter plugin/adapter |
| Assembly | Compiler-generated output, intrinsics or audited dependency only | Exception by specialist approval | No handwritten assembly in pilot product code |

## Why Go for the backend

- A small team can deploy one statically compiled application/worker artifact with predictable operations.
- Concurrency primitives suit bounded HTTP, webhook, worker and provider-I/O workloads, while database transactions remain the correctness boundary.
- Fast builds, standard tooling, explicit errors and straightforward profiling support frequent safe releases.
- Go does not make concurrency automatically safe: goroutine creation, connection pools, queues and provider calls still require hard limits, cancellation, idempotency and load tests.

## Why Flutter/Dart for mobile

- One primary team can deliver Android first and later reuse substantial UI/application infrastructure for iOS.
- Flutter supports Android/iOS and desktop targets and can call native code, but support does not prove Ubi meets D1 device, accessibility, binary-size, background and offline targets.
- The local database/outbox/repository design is framework-independent enough to permit reversal if the prototype fails.
- Customer public web is not assumed to be Flutter web: its server-rendering, first-load, semantic HTML, SEO/link and constrained-network needs favor a document/web architecture.

## C and C++ policy

C and C++ are valuable for mature native libraries, codecs, image/audio/video processing, compression, cryptography implementations, database engines and very tight CPU/memory work. They are not a shortcut to a fast product.

Ubi may introduce C/C++ only when:

1. a representative benchmark identifies a material user/cost problem;
2. platform/framework optimization is insufficient;
3. an audited maintained library is preferable to custom code;
4. supported ABIs/platforms/toolchains are explicit;
5. memory ownership, thread rules, cancellation, input bounds and error mapping are documented;
6. fuzzing, sanitizers, static analysis, dependency/SBOM and malicious-input tests pass;
7. crash symbolication and native performance telemetry work;
8. an owner can patch CVEs and reproduce builds across supported platforms.

Shared native code exposes a versioned, minimal C ABI with fixed-width types, pointer/length pairs, explicit alloc/free ownership and stable error codes. C++ exceptions, RTTI-dependent types, STL containers and allocator ownership do not cross the ABI. A safe Dart/Kotlin/Swift/Rust wrapper validates all inputs and lifetimes.

Booking state, authorization, payment truth, tenant isolation, policy evaluation and offline command semantics remain in Go/Dart/SQL domain contracts, never only in native code.

## Rust policy

Rust is the preferred candidate for **new** memory-sensitive parsers, media helpers, cryptographic glue or cross-platform native components when the team can operate its toolchain. Safe Rust should dominate; every `unsafe` block and FFI boundary needs an invariant comment, review and test. Rust does not eliminate logic, supply-chain or FFI vulnerabilities.

Rust is not required for the initial booking backend simply because it is fast. Introducing both Go and Rust server stacks before a measured need would fragment ownership.

## Java, Kotlin and Swift policy

- Android platform/plugin code uses Kotlin by default; Java is accepted for provider SDKs, legacy APIs or libraries and may be wrapped rather than rewritten.
- iOS/macOS platform/plugin code uses Swift; it does not duplicate the full Flutter application.
- Platform adapters expose narrow capabilities such as secure storage, notifications, background scheduling, links, camera/files and provider SDK calls.
- Every bridge is contract-tested for cancellation, lifecycle, thread, process death, permission denial and error mapping.

## Assembly policy

No handwritten assembly is allowed in first release. Compilers already generate architecture-specific code; vetted native libraries may contain assembly/SIMD. A future exception requires a reproducible benchmark, supported CPU dispatch/fallback, specialist review, constant-time/security analysis where relevant, per-architecture tests and a clear removal path.

## Alternatives rejected for the pilot

- **All-native clients:** maximum platform control but duplicates scarce engineering effort before product validation.
- **Flutter for every web surface:** possible, but not the best default for account-free, server-rendered, semantic, low-JavaScript customer links.
- **C++ or Rust primary backend:** technically capable, but higher early staffing/operational complexity without evidence that Go cannot meet targets.
- **Java primary backend:** capable and mature, but Go offers a smaller operational baseline for the proposed team; Java remains useful at Android/provider boundaries.
- **Many microservices in different languages:** maximizes coordination and failure modes before scale/team boundaries exist.

## Technical baselines checked

Checked 2026-09-18; reverify supported versions/toolchains before implementation:

- Flutter's official [supported-platform matrix](https://docs.flutter.dev/reference/supported-platforms) covers mobile, desktop and web targets, while its [web guidance](https://docs.flutter.dev/platform-integration/web) explicitly notes that not every document-centric HTML scenario is ideally suited to Flutter.
- Android's official [NDK guidance](https://developer.android.com/ndk/guides) positions C/C++ for existing native libraries or genuinely computational/low-latency work—not as the default for ordinary Android application code.
- Dart provides official [C interoperability through `dart:ffi`](https://dart.dev/interop/c-interop); that makes a narrow native boundary possible but does not remove memory/lifecycle risk.
- Go's official guidance describes lightweight goroutines/channels and warns that concurrency is not automatically parallelism or safe resource management; Ubi therefore keeps explicit bounds and database correctness controls in addition to choosing Go.
- Rust's official [unsafety reference](https://doc.rust-lang.org/reference/unsafety.html) and [FFI guidance](https://doc.rust-lang.org/nomicon/ffi.html) require unsafe boundaries to be isolated behind safe abstractions.
- Apple documents [Swift with C/C++ interoperability](https://developer.apple.com/swift), supporting thin native adapters rather than a duplicated app.
- React provides official [server rendering APIs](https://react.dev/reference/react-dom/server), consistent with Ubi's server-rendered customer-link requirement.

## Validation required before approval

- Flutter release prototype on purchased D1/D2 devices passes cold/warm start, frame, binary, memory, local database, offline outbox, background and accessibility gates.
- Go vertical-slice prototype passes booking-concurrency, idempotency, PostgreSQL pool, webhook burst, memory/CPU and graceful-shutdown tests.
- React customer route meets server-rendered content, semantic accessibility and W1/N1 payload/performance budgets.
- One native interop spike proves build/sign/symbol/fuzz/sanitizer/FFI lifecycle on Android; it does not need to ship.
- Team ownership, CI duration, dependency patching and hiring feasibility are acceptable.

## Consequences

- Ubi gains room to use C/C++/Rust without making them default complexity.
- Contract generation and versioning are essential across Go, Dart and TypeScript.
- Native modules are exceptional release artifacts with stricter testing/signing/observability.
- Framework reversal remains possible before large UI/plugin investment.
- Language additions after approval require a new ADR with measured benefit and full lifecycle cost.
