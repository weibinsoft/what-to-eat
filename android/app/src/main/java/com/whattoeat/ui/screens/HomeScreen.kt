package com.whattoeat.ui.screens

import android.media.AudioManager
import android.media.ToneGenerator
import androidx.compose.animation.core.*
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.whattoeat.data.api.models.DecisionRecord
import com.whattoeat.data.api.models.Menu
import com.whattoeat.ui.theme.GradientEnd
import com.whattoeat.ui.theme.GradientStart
import com.whattoeat.viewmodel.HomeViewModel
import kotlinx.coroutines.delay
import java.text.SimpleDateFormat
import java.util.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(
    viewModel: HomeViewModel,
    onLogout: () -> Unit,
    onNavigateToSettings: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()
    var selectedTab by rememberSaveable { mutableStateOf(0) }
    var showAddMenuDialog by remember { mutableStateOf(false) }
    var showDecisionResultDialog by remember { mutableStateOf(false) }

    // 音效生成器
    val toneGenerator = remember {
        try {
            ToneGenerator(AudioManager.STREAM_MUSIC, 80)
        } catch (e: Exception) {
            null
        }
    }

    // 清理音效资源
    DisposableEffect(Unit) {
        onDispose {
            toneGenerator?.release()
        }
    }

    // 决策过程中播放连续的滚动音效（类似老虎机）
    LaunchedEffect(uiState.isDeciding) {
        if (uiState.isDeciding) {
            var tickCount = 0
            while (uiState.isDeciding && tickCount < 40) {
                try {
                    // 使用 DTMF 音调模拟老虎机滚动声
                    val tones = listOf(
                        ToneGenerator.TONE_DTMF_1,
                        ToneGenerator.TONE_DTMF_2,
                        ToneGenerator.TONE_DTMF_3,
                        ToneGenerator.TONE_DTMF_4,
                        ToneGenerator.TONE_DTMF_5
                    )
                    toneGenerator?.startTone(tones[tickCount % tones.size], 50)
                } catch (e: Exception) {
                    // 忽略音效播放错误
                }
                // 逐渐减速
                val delayMs = 80L + (tickCount * 5L)
                delay(delayMs)
                tickCount++
            }
        }
    }

    // 监听决策结果并播放成功音效
    LaunchedEffect(uiState.decisionResult) {
        if (uiState.decisionResult != null) {
            // 播放成功音效（庆祝音）
            try {
                toneGenerator?.startTone(ToneGenerator.TONE_PROP_BEEP2, 300)
            } catch (e: Exception) {
                // 忽略音效播放错误
            }
            showDecisionResultDialog = true
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        "今天吃什么",
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = GradientStart,
                    titleContentColor = Color.White
                ),
                actions = {
                    IconButton(onClick = onNavigateToSettings) {
                        Icon(
                            Icons.Default.Settings,
                            contentDescription = "设置",
                            tint = Color.White
                        )
                    }
                    IconButton(onClick = onLogout) {
                        Icon(
                            Icons.Default.Logout,
                            contentDescription = "退出",
                            tint = Color.White
                        )
                    }
                }
            )
        },
        floatingActionButton = {
            if (selectedTab == 1) {
                FloatingActionButton(
                    onClick = { showAddMenuDialog = true },
                    containerColor = GradientStart
                ) {
                    Icon(
                        Icons.Default.Add,
                        contentDescription = "添加菜单",
                        tint = Color.White
                    )
                }
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            // Tab 栏
            TabRow(
                selectedTabIndex = selectedTab,
                containerColor = MaterialTheme.colorScheme.surface,
                contentColor = GradientStart
            ) {
                Tab(
                    selected = selectedTab == 0,
                    onClick = { selectedTab = 0 },
                    text = { Text("决策") },
                    icon = { Icon(Icons.Default.Casino, contentDescription = null) }
                )
                Tab(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    text = { Text("菜单") },
                    icon = { Icon(Icons.Default.Restaurant, contentDescription = null) }
                )
                Tab(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    text = { Text("历史") },
                    icon = { Icon(Icons.Default.History, contentDescription = null) }
                )
            }

            // 内容区域
            when (selectedTab) {
                0 -> DecisionTab(
                    slotDisplayText = uiState.slotDisplayText,
                    isDeciding = uiState.isDeciding,
                    menuCount = uiState.menus.size,
                    onDecide = { viewModel.decide() }
                )
                1 -> MenuListTab(
                    menus = uiState.menus,
                    isLoading = uiState.isLoading,
                    onDelete = { viewModel.deleteMenu(it) },
                    onRefresh = { viewModel.loadData() }
                )
                2 -> HistoryTab(
                    records = uiState.historyRecords,
                    isLoading = uiState.isLoading,
                    onRefresh = { viewModel.loadData() }
                )
            }
        }
    }

    // 添加菜单对话框
    if (showAddMenuDialog) {
        AddMenuDialog(
            restaurants = uiState.restaurants.map { it.name },
            isLoading = uiState.isAddingMenu,
            onDismiss = { showAddMenuDialog = false },
            onConfirm = { restaurantName, dishName ->
                viewModel.addMenu(restaurantName, dishName)
                showAddMenuDialog = false
            }
        )
    }

    // 决策结果对话框
    if (showDecisionResultDialog && uiState.decisionResult != null) {
        DecisionResultDialog(
            result = uiState.decisionResult!!,
            message = uiState.decisionMessage,
            onDismiss = {
                showDecisionResultDialog = false
                viewModel.clearDecisionResult()
            }
        )
    }

    // 错误提示
    uiState.error?.let { error ->
        LaunchedEffect(error) {
            // 显示 Snackbar
        }
    }
}

@Composable
fun DecisionTab(
    slotDisplayText: String,
    isDeciding: Boolean,
    menuCount: Int,
    onDecide: () -> Unit
) {
    val infiniteTransition = rememberInfiniteTransition(label = "slot")
    val scale by infiniteTransition.animateFloat(
        initialValue = 1f,
        targetValue = if (isDeciding) 1.1f else 1f,
        animationSpec = infiniteRepeatable(
            animation = tween(200),
            repeatMode = RepeatMode.Reverse
        ),
        label = "scale"
    )

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(
                brush = Brush.verticalGradient(
                    colors = listOf(
                        MaterialTheme.colorScheme.surface,
                        MaterialTheme.colorScheme.surface.copy(alpha = 0.95f)
                    )
                )
            )
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        // 老虎机显示区域
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .height(200.dp)
                .scale(if (isDeciding) scale else 1f),
            shape = RoundedCornerShape(20.dp),
            colors = CardDefaults.cardColors(
                containerColor = Color.White
            ),
            elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
        ) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(
                        brush = Brush.verticalGradient(
                            colors = listOf(GradientStart.copy(alpha = 0.1f), GradientEnd.copy(alpha = 0.1f))
                        )
                    ),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = slotDisplayText,
                    fontSize = if (slotDisplayText.length > 10) 20.sp else 28.sp,
                    fontWeight = FontWeight.Bold,
                    textAlign = TextAlign.Center,
                    modifier = Modifier.padding(16.dp)
                )
            }
        }

        Spacer(modifier = Modifier.height(32.dp))

        // 菜单数量提示
        Text(
            text = "共 $menuCount 道菜可选",
            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
        )

        Spacer(modifier = Modifier.height(24.dp))

        // 决策按钮
        Button(
            onClick = onDecide,
            modifier = Modifier
                .fillMaxWidth()
                .height(60.dp),
            enabled = !isDeciding && menuCount > 0,
            shape = RoundedCornerShape(16.dp),
            colors = ButtonDefaults.buttonColors(
                containerColor = GradientStart
            )
        ) {
            if (isDeciding) {
                Text("🎰 选择中...", fontSize = 18.sp)
            } else {
                Icon(Icons.Default.Casino, contentDescription = null)
                Spacer(modifier = Modifier.width(8.dp))
                Text("开始决策", fontSize = 18.sp)
            }
        }
    }
}

@Composable
fun MenuListTab(
    menus: List<Menu>,
    isLoading: Boolean,
    onDelete: (Long) -> Unit,
    onRefresh: () -> Unit
) {
    if (isLoading && menus.isEmpty()) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Text("加载中...", color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f))
        }
    } else if (menus.isEmpty()) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Text(
                    text = "🍽️",
                    fontSize = 64.sp
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "还没有菜单",
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                )
                Text(
                    text = "点击右下角添加",
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.4f),
                    fontSize = 14.sp
                )
            }
        }
    } else {
        // 按餐厅分组
        val groupedMenus = menus.groupBy { it.restaurant?.name ?: "未知餐厅" }

        LazyColumn(
            modifier = Modifier.fillMaxSize(),
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            groupedMenus.forEach { (restaurantName, restaurantMenus) ->
                item {
                    Text(
                        text = restaurantName,
                        fontWeight = FontWeight.Bold,
                        fontSize = 16.sp,
                        color = GradientStart,
                        modifier = Modifier.padding(vertical = 8.dp)
                    )
                }
                items(restaurantMenus) { menu ->
                    MenuItemCard(
                        menu = menu,
                        onDelete = { onDelete(menu.id) }
                    )
                }
            }
        }
    }
}

@Composable
fun MenuItemCard(
    menu: Menu,
    onDelete: () -> Unit
) {
    var showDeleteConfirm by remember { mutableStateOf(false) }

    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = menu.dishName,
                fontSize = 16.sp
            )
            IconButton(
                onClick = { showDeleteConfirm = true }
            ) {
                Icon(
                    Icons.Default.Delete,
                    contentDescription = "删除",
                    tint = MaterialTheme.colorScheme.error
                )
            }
        }
    }

    if (showDeleteConfirm) {
        AlertDialog(
            onDismissRequest = { showDeleteConfirm = false },
            title = { Text("确认删除") },
            text = { Text("确定要删除 \"${menu.dishName}\" 吗？") },
            confirmButton = {
                TextButton(
                    onClick = {
                        onDelete()
                        showDeleteConfirm = false
                    }
                ) {
                    Text("删除", color = MaterialTheme.colorScheme.error)
                }
            },
            dismissButton = {
                TextButton(onClick = { showDeleteConfirm = false }) {
                    Text("取消")
                }
            }
        )
    }
}

@Composable
fun HistoryTab(
    records: List<DecisionRecord>,
    isLoading: Boolean,
    onRefresh: () -> Unit
) {
    if (isLoading && records.isEmpty()) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Text("加载中...", color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f))
        }
    } else if (records.isEmpty()) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Text(
                    text = "📝",
                    fontSize = 64.sp
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "还没有决策记录",
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                )
            }
        }
    } else {
        LazyColumn(
            modifier = Modifier.fillMaxSize(),
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            items(records) { record ->
                HistoryItemCard(record = record)
            }
        }
    }
}

@Composable
fun HistoryItemCard(record: DecisionRecord) {
    val dateFormat = remember { SimpleDateFormat("MM-dd HH:mm", Locale.getDefault()) }
    val formattedDate = try {
        val inputFormat = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss", Locale.getDefault())
        val date = inputFormat.parse(record.decidedAt.substring(0, 19))
        date?.let { dateFormat.format(it) } ?: record.decidedAt
    } catch (e: Exception) {
        record.decidedAt
    }

    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column {
                Text(
                    text = record.menu?.dishName ?: "未知菜品",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Medium
                )
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = record.menu?.restaurant?.name ?: "未知餐厅",
                    fontSize = 14.sp,
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                )
            }
            Text(
                text = formattedDate,
                fontSize = 12.sp,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
            )
        }
    }
}

@OptIn(ExperimentalLayoutApi::class)
@Composable
fun AddMenuDialog(
    restaurants: List<String>,
    isLoading: Boolean,
    onDismiss: () -> Unit,
    onConfirm: (String, String) -> Unit
) {
    var restaurantName by remember { mutableStateOf("") }
    var dishName by remember { mutableStateOf("") }

    // 过滤匹配的餐厅
    val filteredRestaurants = remember(restaurantName, restaurants) {
        if (restaurantName.isBlank()) {
            restaurants.take(5) // 显示前5个历史餐厅
        } else {
            restaurants.filter { it.contains(restaurantName, ignoreCase = true) }.take(5)
        }
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("添加菜单") },
        text = {
            Column(
                modifier = Modifier.verticalScroll(rememberScrollState())
            ) {
                // 餐厅输入
                OutlinedTextField(
                    value = restaurantName,
                    onValueChange = { restaurantName = it },
                    label = { Text("餐厅名称") },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true
                )

                // 历史餐厅快捷选择
                if (filteredRestaurants.isNotEmpty() && restaurantName.isEmpty()) {
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(
                        text = "历史餐厅（点击选择）",
                        fontSize = 12.sp,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                    )
                    Spacer(modifier = Modifier.height(4.dp))
                    FlowRow(
                        horizontalArrangement = Arrangement.spacedBy(8.dp),
                        verticalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        filteredRestaurants.forEach { restaurant ->
                            OutlinedButton(
                                onClick = { restaurantName = restaurant },
                                contentPadding = PaddingValues(horizontal = 12.dp, vertical = 4.dp),
                                modifier = Modifier.height(32.dp)
                            ) {
                                Text(restaurant, fontSize = 12.sp)
                            }
                        }
                    }
                }

                Spacer(modifier = Modifier.height(16.dp))

                // 菜品输入
                OutlinedTextField(
                    value = dishName,
                    onValueChange = { dishName = it },
                    label = { Text("菜品名称") },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true
                )
            }
        },
        confirmButton = {
            Button(
                onClick = { onConfirm(restaurantName, dishName) },
                enabled = restaurantName.isNotBlank() && dishName.isNotBlank() && !isLoading
            ) {
                Text(if (isLoading) "添加中..." else "添加")
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text("取消")
            }
        }
    )
}

@Composable
fun DecisionResultDialog(
    result: String,
    message: String?,
    onDismiss: () -> Unit
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Text(
                text = "🎉 今天吃这个！",
                textAlign = TextAlign.Center,
                modifier = Modifier.fillMaxWidth()
            )
        },
        text = {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(
                    text = result,
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    textAlign = TextAlign.Center
                )
                if (message != null) {
                    Spacer(modifier = Modifier.height(16.dp))
                    Text(
                        text = message,
                        fontSize = 14.sp,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f),
                        textAlign = TextAlign.Center
                    )
                }
            }
        },
        confirmButton = {
            Button(onClick = onDismiss) {
                Text("好的")
            }
        }
    )
}
