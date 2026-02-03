package com.whattoeat.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.whattoeat.data.api.models.DecisionRecord
import com.whattoeat.data.api.models.Menu
import com.whattoeat.data.api.models.Restaurant
import com.whattoeat.data.repository.Result
import com.whattoeat.data.repository.WhatToEatRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject
import kotlin.random.Random

data class HomeUiState(
    val isLoading: Boolean = false,
    val menus: List<Menu> = emptyList(),
    val restaurants: List<Restaurant> = emptyList(),
    val historyRecords: List<DecisionRecord> = emptyList(),
    val error: String? = null,
    // 决策相关
    val isDeciding: Boolean = false,
    val decisionResult: String? = null,
    val decisionMessage: String? = null,
    val slotDisplayText: String = "🍜",
    // 添加菜单
    val isAddingMenu: Boolean = false,
    val addMenuSuccess: Boolean = false
)

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val repository: WhatToEatRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(HomeUiState())
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()

    private val foodEmojis = listOf("🍜", "🍕", "🍔", "🍣", "🍱", "🍛", "🍝", "🍲", "🥘", "🥡", "🍙", "🍚")

    init {
        loadData()
    }

    fun loadData() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)

            // 并行加载数据
            val menusResult = repository.getMenus()
            val restaurantsResult = repository.getRestaurants()
            val historyResult = repository.getHistory()

            _uiState.value = _uiState.value.copy(
                isLoading = false,
                menus = when (menusResult) {
                    is Result.Success -> menusResult.data
                    is Result.Error -> emptyList()
                },
                restaurants = when (restaurantsResult) {
                    is Result.Success -> restaurantsResult.data
                    is Result.Error -> emptyList()
                },
                historyRecords = when (historyResult) {
                    is Result.Success -> historyResult.data.records
                    is Result.Error -> emptyList()
                }
            )
        }
    }

    fun decide() {
        if (_uiState.value.menus.isEmpty()) {
            _uiState.value = _uiState.value.copy(error = "请先添加菜单")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(
                isDeciding = true,
                decisionResult = null,
                decisionMessage = null
            )

            // 老虎机动画效果
            val menuDisplayNames = _uiState.value.menus.map { menu ->
                "${menu.restaurant?.name ?: "未知"} - ${menu.dishName}"
            }

            repeat(20) { i ->
                val randomIndex = Random.nextInt(menuDisplayNames.size)
                _uiState.value = _uiState.value.copy(
                    slotDisplayText = menuDisplayNames[randomIndex]
                )
                delay(50L + i * 10L) // 逐渐减速
            }

            // 调用 API 获取真正的结果
            when (val result = repository.decide()) {
                is Result.Success -> {
                    val menu = result.data.menu
                    val displayText = "${menu.restaurant?.name ?: "未知"} - ${menu.dishName}"
                    _uiState.value = _uiState.value.copy(
                        isDeciding = false,
                        slotDisplayText = displayText,
                        decisionResult = displayText,
                        decisionMessage = result.data.message
                    )
                    // 刷新历史记录
                    loadHistory()
                }
                is Result.Error -> {
                    _uiState.value = _uiState.value.copy(
                        isDeciding = false,
                        error = result.message
                    )
                }
            }
        }
    }

    fun addMenu(restaurantName: String, dishName: String) {
        if (restaurantName.isBlank() || dishName.isBlank()) {
            _uiState.value = _uiState.value.copy(error = "请输入餐厅名和菜品名")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isAddingMenu = true)

            when (val result = repository.createMenu(restaurantName, dishName)) {
                is Result.Success -> {
                    _uiState.value = _uiState.value.copy(
                        isAddingMenu = false,
                        addMenuSuccess = true
                    )
                    loadData() // 重新加载数据
                }
                is Result.Error -> {
                    _uiState.value = _uiState.value.copy(
                        isAddingMenu = false,
                        error = result.message
                    )
                }
            }
        }
    }

    fun deleteMenu(menuId: Long) {
        viewModelScope.launch {
            when (val result = repository.deleteMenu(menuId)) {
                is Result.Success -> {
                    loadData() // 重新加载数据
                }
                is Result.Error -> {
                    _uiState.value = _uiState.value.copy(error = result.message)
                }
            }
        }
    }

    private fun loadHistory() {
        viewModelScope.launch {
            when (val result = repository.getHistory()) {
                is Result.Success -> {
                    _uiState.value = _uiState.value.copy(
                        historyRecords = result.data.records
                    )
                }
                is Result.Error -> { /* 忽略历史加载错误 */ }
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }

    fun clearAddMenuSuccess() {
        _uiState.value = _uiState.value.copy(addMenuSuccess = false)
    }

    fun clearDecisionResult() {
        _uiState.value = _uiState.value.copy(
            decisionResult = null,
            decisionMessage = null,
            slotDisplayText = foodEmojis.random()
        )
    }
}
