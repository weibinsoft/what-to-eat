package com.whattoeat.data.repository

import com.whattoeat.data.api.ApiService
import com.whattoeat.data.api.models.*
import com.whattoeat.data.datastore.SettingsDataStore
import java.net.ConnectException
import java.net.SocketTimeoutException
import java.net.UnknownHostException
import javax.inject.Inject
import javax.inject.Singleton

sealed class Result<out T> {
    data class Success<T>(val data: T) : Result<T>()
    data class Error(val message: String) : Result<Nothing>()
}

@Singleton
class WhatToEatRepository @Inject constructor(
    private val apiService: ApiService,
    private val settingsDataStore: SettingsDataStore
) {

    private fun getErrorMessage(e: Exception): String {
        return when (e) {
            is UnknownHostException -> "无法连接到服务器，请检查服务器地址"
            is ConnectException -> "连接服务器失败，请检查网络或服务器地址"
            is SocketTimeoutException -> "连接超时，请检查网络"
            else -> e.message ?: "网络错误"
        }
    }

    // 认证
    suspend fun register(username: String, password: String): Result<User> {
        return try {
            val response = apiService.register(RegisterRequest(username, password))
            if (response.isSuccessful && response.body()?.code == 0) {
                response.body()?.data?.let { user ->
                    Result.Success(user)
                } ?: Result.Error("注册失败")
            } else {
                Result.Error(response.body()?.message ?: "注册失败")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }

    suspend fun login(username: String, password: String): Result<LoginResponse> {
        return try {
            val response = apiService.login(LoginRequest(username, password))
            if (response.isSuccessful && response.body()?.code == 0) {
                response.body()?.data?.let { loginResponse ->
                    // 保存登录信息
                    settingsDataStore.setToken(loginResponse.token)
                    settingsDataStore.setUserInfo(loginResponse.userId, loginResponse.username)
                    Result.Success(loginResponse)
                } ?: Result.Error("登录失败")
            } else {
                Result.Error(response.body()?.message ?: "用户名或密码错误")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }

    suspend fun logout() {
        settingsDataStore.clearUserInfo()
    }

    // 菜单
    suspend fun getMenus(): Result<List<Menu>> {
        return try {
            val response = apiService.getMenus()
            if (response.isSuccessful && response.body()?.code == 0) {
                Result.Success(response.body()?.data ?: emptyList())
            } else {
                Result.Error(response.body()?.message ?: "获取菜单失败")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }

    suspend fun createMenu(restaurantName: String, dishName: String): Result<CreateMenuResponse> {
        return try {
            val response = apiService.createMenu(CreateMenuRequest(restaurantName, dishName))
            if (response.isSuccessful && response.body()?.code == 0) {
                response.body()?.data?.let { data ->
                    Result.Success(data)
                } ?: Result.Error("添加失败")
            } else {
                Result.Error(response.body()?.message ?: "添加失败")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }

    suspend fun deleteMenu(id: Long): Result<Unit> {
        return try {
            val response = apiService.deleteMenu(id)
            if (response.isSuccessful && response.body()?.code == 0) {
                Result.Success(Unit)
            } else {
                Result.Error(response.body()?.message ?: "删除失败")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }

    // 餐厅
    suspend fun getRestaurants(): Result<List<Restaurant>> {
        return try {
            val response = apiService.getRestaurants()
            if (response.isSuccessful && response.body()?.code == 0) {
                Result.Success(response.body()?.data ?: emptyList())
            } else {
                Result.Error(response.body()?.message ?: "获取餐厅失败")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }

    // 决策
    suspend fun decide(): Result<DecideResponse> {
        return try {
            val response = apiService.decide(DecideRequest())
            if (response.isSuccessful && response.body()?.code == 0) {
                response.body()?.data?.let { data ->
                    Result.Success(data)
                } ?: Result.Error("决策失败")
            } else {
                Result.Error(response.body()?.message ?: "决策失败")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }

    suspend fun getHistory(): Result<HistoryResponse> {
        return try {
            val response = apiService.getHistory()
            if (response.isSuccessful && response.body()?.code == 0) {
                response.body()?.data?.let { data ->
                    Result.Success(data)
                } ?: Result.Success(HistoryResponse(emptyList(), 0))
            } else {
                Result.Error(response.body()?.message ?: "获取历史失败")
            }
        } catch (e: Exception) {
            Result.Error(getErrorMessage(e))
        }
    }
}
