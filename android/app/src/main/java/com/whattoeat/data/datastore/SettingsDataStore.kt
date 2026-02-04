package com.whattoeat.data.datastore

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.longPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "settings")

@Singleton
class SettingsDataStore @Inject constructor(
    @ApplicationContext private val context: Context
) {
    companion object {
        private val KEY_SERVER_HOST = stringPreferencesKey("server_host")
        private val KEY_TOKEN = stringPreferencesKey("token")
        private val KEY_USER_ID = longPreferencesKey("user_id")
        private val KEY_USERNAME = stringPreferencesKey("username")

        const val DEFAULT_SERVER_HOST = "http://10.0.2.2:8080" // Android 模拟器 localhost

        // 游客用户信息
        const val GUEST_USER_ID = 1L
        const val GUEST_USERNAME = "guest"
    }

    // 服务器地址
    val serverHost: Flow<String> = context.dataStore.data.map { preferences ->
        preferences[KEY_SERVER_HOST] ?: DEFAULT_SERVER_HOST
    }

    suspend fun getServerHost(): String {
        return context.dataStore.data.first()[KEY_SERVER_HOST] ?: DEFAULT_SERVER_HOST
    }

    suspend fun setServerHost(host: String) {
        context.dataStore.edit { preferences ->
            preferences[KEY_SERVER_HOST] = host
        }
    }

    // Token
    val token: Flow<String?> = context.dataStore.data.map { preferences ->
        preferences[KEY_TOKEN]
    }

    suspend fun getToken(): String? {
        return context.dataStore.data.first()[KEY_TOKEN]
    }

    suspend fun setToken(token: String?) {
        context.dataStore.edit { preferences ->
            if (token != null) {
                preferences[KEY_TOKEN] = token
            } else {
                preferences.remove(KEY_TOKEN)
            }
        }
    }

    // 用户信息
    val userId: Flow<Long?> = context.dataStore.data.map { preferences ->
        preferences[KEY_USER_ID]
    }

    val username: Flow<String?> = context.dataStore.data.map { preferences ->
        preferences[KEY_USERNAME]
    }

    suspend fun setUserInfo(userId: Long, username: String) {
        context.dataStore.edit { preferences ->
            preferences[KEY_USER_ID] = userId
            preferences[KEY_USERNAME] = username
        }
    }

    suspend fun clearUserInfo() {
        context.dataStore.edit { preferences ->
            preferences.remove(KEY_TOKEN)
            preferences.remove(KEY_USER_ID)
            preferences.remove(KEY_USERNAME)
        }
    }

    // 检查是否已登录
    suspend fun isLoggedIn(): Boolean {
        return getToken() != null
    }
}
