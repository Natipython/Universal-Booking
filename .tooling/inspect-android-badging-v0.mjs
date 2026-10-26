// aapt2 uses minSdkVersion, not sdkVersion. Parse the exact labelled fields.
export function inspectApkBadging(text) {
  const packageLine = text.match(/^package:.*$/m)?.[0];
  const packageName = packageLine?.match(/\bname='([^']+)'/)?.[1];
  const minSdk = Number(text.match(/^minSdkVersion:'(\d+)'$/m)?.[1]);
  const targetSdk = Number(text.match(/^targetSdkVersion:'(\d+)'$/m)?.[1]);
  if (!packageName || !Number.isInteger(minSdk) || !Number.isInteger(targetSdk)) {
    throw Error('APK badging lacks package/minimum/target API fields');
  }
  return { package: packageLine, packageName, minSdk, targetSdk,
    nativeCode: text.match(/^native-code:.*$/m)?.[0],
    debuggable: /^application-debuggable$/m.test(text),
    permissions: [...text.matchAll(/^uses-permission: name='([^']+)'/gm)].map(match => match[1]) };
}
