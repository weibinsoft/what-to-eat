# Add project specific ProGuard rules here.
# By default, the flags in this file are appended to flags specified
# in /sdk/tools/proguard/proguard-android.txt

# Retrofit
-keepattributes Signature
-keepattributes Exceptions
-keep class com.whattoeat.data.api.models.** { *; }

# Gson
-keep class com.google.gson.** { *; }
-keepattributes *Annotation*
