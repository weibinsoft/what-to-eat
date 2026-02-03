package com.whattoeat.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.whattoeat.data.datastore.SettingsDataStore
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
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
            _uiState.value = _uiState.value.copy(isSaving = true)
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
}
