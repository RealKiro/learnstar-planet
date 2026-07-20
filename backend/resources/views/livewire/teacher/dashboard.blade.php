<div>
    <!-- 班级选择器 -->
    <div class="flex items-center gap-4 mb-8">
        <flux:heading size="xl" class="bc-gradient-text">数据看板</flux:heading>
        <span class="free-badge">全免费</span>
        
        <flux:select wire:model.live="classRoomId" class="ml-auto w-48">
            @foreach($classRooms as $class)
                <flux:select.option value="{{ $class->id }}">{{ $class->name }}</flux:select.option>
            @endforeach
        </flux:select>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <flux:card class="bc-card bc-magnetic">
            <flux:text.class="text-sm opacity-60">总积分</flux:text>
            <flux:heading size="lg">{{ $stats['net_total'] ?? 0 }}</flux:heading>
            <flux:text class="text-xs opacity-40">加分 {{ $stats['total_add'] ?? 0 }} / 减分 {{ $stats['total_subtract'] ?? 0 }}</flux:text>
        </flux:card>
        
        <flux:card class="bc-card bc-magnetic">
            <flux:text class="text-sm opacity-60">今日加分</flux:text>
            <flux:heading size="lg" class="text-green-600">{{ $stats['today_add'] ?? 0 }}</flux:heading>
            <flux:text class="text-xs opacity-40">今日新增积分</flux:text>
        </flux:card>
        
        <flux:card class="bc-card bc-magnetic">
            <flux:text class="text-sm opacity-60">积分记录</flux:text>
            <flux:heading size="lg">{{ $stats['score_count'] ?? 0 }}</flux:heading>
            <flux:text class="text-xs opacity-40">累计操作次数</flux:text>
        </flux:card>
        
        <flux:card class="bc-card bc-magnetic">
            <flux:text class="text-sm opacity-60">宠物数量</flux:text>
            <flux:heading size="lg">{{ array_sum($petDistribution) }}</flux:heading>
            <flux:text class="text-xs opacity-40">全班宠物总数</flux:text>
        </flux:card>
    </div>

    <!-- 两栏布局：排行榜 + 积分动态 -->
    <div class="grid md:grid-cols-2 gap-6">
        <!-- TOP 10 排行榜 -->
        <flux:card class="bc-card">
            <flux:heading size="md" class="mb-4">⭐ TOP 10 排行榜</flux:heading>
            
            @foreach($topStudents as $index => $student)
                <div class="flex items-center gap-3 py-2 {{ $index > 0 ? 'border-t' : '' }}">
                    <span class="w-8 text-center font-bold {{ $index < 3 ? 'bc-gradient-text' : 'opacity-60' }}">
                        {{ $student['rank'] }}
                    </span>
                    <span class="text-xl">{{ $student['pet_emoji'] }}</span>
                    <span class="flex-1 font-medium">{{ $student['name'] }}</span>
                    <span class="font-bold {{ $student['total_score'] > 0 ? 'text-green-600' : 'text-red-500' }}">
                        {{ $student['total_score'] }}分
                    </span>
                </div>
            @endforeach
            
            @if($topStudents->isEmpty())
                <flux:text class="text-center opacity-40 py-8">暂无排行数据</flux:text>
            @endif
        </flux:card>

        <!-- 最近积分动态 -->
        <flux:card class="bc-card">
            <flux:heading size="md" class="mb-4">📊 积分动态</flux:heading>
            
            @foreach($stats['recent_scores'] ?? [] as $score)
                <div class="flex items-center gap-3 py-2 border-t">
                    <span class="{{ $score->type === 'add' ? 'text-green-600' : 'text-red-500' }} font-bold">
                        {{ $score->type === 'add' ? '+' : '-' }}{{ $score->amount }}
                    </span>
                    <span class="flex-1">{{ $score->student?->name }}</span>
                    <span class="text-xs opacity-40">{{ $score->rule_code }}</span>
                    <span class="text-xs opacity-40">{{ $score->created_at->format('H:i') }}</span>
                </div>
            @endforeach
        </flux:card>
    </div>

    <!-- 积分趋势图区域 -->
    <flux:card class="bc-card mt-6">
        <flux:heading size="md" class="mb-4">📈 近7天积分趋势</flux:heading>
        <div class="h-48 flex items-end gap-2">
            @foreach($scoreTrend as $day)
                <div class="flex-1 flex flex-col items-center gap-1">
                    <div class="w-full bg-gradient-to-t from-indigo-500 to-indigo-300 rounded-t-md"
                         style="height: {{ max(4, ($day['add'] / max(1, collect($scoreTrend)->max('add'))) * 160) }}px">
                    </div>
                    <span class="text-xs opacity-40">{{ \Carbon\Carbon::parse($day['date'])->format('M/d') }}</span>
                </div>
            @endforeach
        </div>
    </flux:card>
</div>
