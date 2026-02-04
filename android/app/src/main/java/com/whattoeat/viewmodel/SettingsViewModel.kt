package com.whattoeat.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.whattoeat.data.api.ApiService
import com.whattoeat.data.datastore.SettingsDataStore
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.net.ConnectException
import java.net.SocketTimeoutException
import java.net.UnknownHostException
import java.util.concurrent.TimeUnit
import javax.inject.Inject

data class SettingsUiState(
    val serverHost: String = "",
    val isSaving: Boolean = false,
    val saveSuccess: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class SettingsViewModel @Inject constructor(
    private val settingsDataStore: SettingsDataStore
) : ViewModel() {

    private val _uiState = MutableStateFlow(SettingsUiState())
    val uiState: StateFlow<SettingsUiState> = _uiState.asStateFlow()

    init {
        loadSettings()
    }

    private fun loadSettings() {
        viewModelScope.launch {
            val host = settingsDataStore.getServerHost()
            _uiState.value = _uiState.value.copy(serverHost = host)
        }
    }

    fun updateServerHost(host: String) {
        _uiState.value = _uiState.value.copy(serverHost = host)
    }

    fun saveSettings() {
        val host = _uiState.value.serverHost.trim()

        if (host.isBlank()) {
            _uiState.value = _uiState.value.copy(error = "请输入服务器地址")
            return
        }

        // 简单验证 URL 格式
        if (!host.startsWith("http://") && !host.startsWith("https://")) {
            _uiState.value = _uiState.value.copy(error = "请输入有效的服务器地址（以 http:// 或 https:// 开头）")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isSaving = true, error = null)
            
            // 先进行健康检查
            val healthCheckResult = checkServerHealth(host)
            
            if (healthCheckResult.isSuccess) {
                // 健康检查通过，保存设置
                try {
                    settingsDataStore.setServerHost(host)
                    _uiState.value = _uiState.value.copy(
                        isSaving = false,
                        saveSuccess = true
                    )
                } catch (e: Exception) {
                    _uiState.value = _uiState.value.copy(
                        isSaving = false,
                        error = "保存失败: ${e.message}"
                    )
                }
            } else {
                // 健康检查失败
                _uiState.value = _uiState.value.copy(
                    isSaving = false,
                    error = healthCheckResult.errorMessage
                )
            }
        }
    }

    private suspend fun checkServerHealth(host: String): HealthCheckResult {
        return withContext(Dispatchers.IO) {
            try {
                // 创建临时的 OkHttpClient
                val client = OkHttpClient.Builder()
                    .connectTimeout(10, TimeUnit.SECONDS)
                    .readTimeout(10, TimeUnit.SECONDS)
                    .build()

                // 创建临时的 Retrofit
                val baseUrl = if (host.endsWith("/")) host else "$host/"
                val retrofit = Retrofit.Builder()
                    .baseUrl(baseUrl)
                    .client(client)
                    .addConverterFactory(GsonConverterFactory.create())
                    .build()

                val apiService = retrofit.create(ApiService::class.java)
                
                // 调用健康检查接口
                val response = apiService.health()
                
                if (response.isSuccessful && response.body()?.status == "ok") {
                    HealthCheckResult(isSuccess = true)
                } else {
                    HealthCheckResult(
                        isSuccess = false,
                        errorMessage = "服务器响应异常，请检查服务器状态"
                    )
                }
            } catch (e: UnknownHostException) {
                HealthCheckResult(
                    isSuccess = false,
                    errorMessage = "无法解析服务器地址，请检查地址是否正确"
                )
            } catch (e: ConnectException) {
                HealthCheckResult(
                    isSuccess = false,
                    errorMessage = "无法连接到服务器，请检查服务器是否启动"
                )
            } catch (e: SocketTimeoutException) {
                HealthCheckResult(
                    isSuccess = false,
                    errorMessage = "连接超时，请检查网络或服务器地址"
                )
            } catch (e: Exception) {
                HealthCheckResult(
                    isSuccess = false,
                    errorMessage = "连接失败: ${e.message ?: "未知错误"}"
                )
            }
        }
    }

    fun resetToDefault() {
        _uiState.value = _uiState.value.copy(
            serverHost = SettingsDataStore.DEFAULT_SERVER_HOST
        )
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }

    fun clearSaveSuccess() {
        _uiState.value = _uiState.value.copy(saveSuccess = false)
    }

    private data class HealthCheckResult(
        val isSuccess: Boolean,
        val errorMessage: String? = null
    )
}
