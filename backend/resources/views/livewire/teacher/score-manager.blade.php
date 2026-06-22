<div>
    <div class="flex items-center justify-between mb-6">
        <flux:heading size="xl" class="bc-gradient-text">积分管理</flux:heading>
        <span class="free-badge">全免费</span>
        <flux:button wire:click="$set('showBatchModal', true)" variant="primary" icon="users">
            批量加分
        </flux:button>
    </div>

    <!-- 快捷加减分面板 -->
    <flux:card class="bc-card mb-6">
        <flux:heading size="md" class="mb-4">⚡ 快捷加分</flux:heading>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
            @foreach($rules as $code => $rule)
                @if($rule['type'] === 'add')
                    <flux:button 
                        wire:click="quickAdd({{ $selectedStudentId ?? 0 }}, '{{ $code }}')"
                        variant="soft"
                        class="bc-magnetic"
                        size="sm">
                        {{ $rule['name'] }} +{{ $rule['amount'] }}
                    </flux:button>
                @endif
            @endforeach
        </div>
    </flux:card>

    <!-- 学生列表 -->
    <flux:card class="bc-card">
        <flux:heading size="md" class="mb-4">📋 学生积分列表</flux:heading>
        
        <div class="space-y-2">
            @foreach($students as $student)
                <div class="flex items-center gap-3 p-3 rounded-xl hover:bg-gray-50 transition"
                     wire:click="$set('selectedStudentId', {{ $student->id }})"
                     class="{{ $selectedStudentId === $student->id ? 'ring-2 ring-indigo-500 bg-indigo-50' : '' }}">
                    
                    <!-- 选择框（批量模式） -->
                    <input type="checkbox" 
                           wire:model="selectedStudents"
                           value="{{ $student->id }}"
                           class="w-4 h-4 rounded">
                    
                    <!-- 宠物等级 -->
                    <span class="text-xl">
                        {{ \App\Models\Student::getEvolutionPath()[$student->pet_level]['emoji'] ?? '🥚' }}
                    </span>
                    
                    <!-- 学生信息 -->
                    <div class="flex-1">
                        <span class="font-medium">{{ $student->name }}</span>
                        <span class="text-xs opacity-40 ml-2">Lv.{{ $student->pet_level }}</span>
                    </div>
                    
                    <!-- 总积分 -->
                    <span class="font-bold {{ $student->total_score > 0 ? 'text-green-600' : 'text-red-500' }}">
                        {{ $student->total_score }}分
                    </span>
                    
                    <!-- 月积分 -->
                    <span class="text-sm opacity-60">
                        本月 {{ $student->monthly_score }}分
                    </span>
                    
                    <!-- 快捷按钮 -->
                    <flux:button.group>
                        <flux:button icon="plus" size="xs" variant="soft"
                            wire:click="quickAdd({{ $student->id }}, 'homework_complete')">
                        </flux:button>
                        <flux:button icon="minus" size="xs" variant="soft" color="danger"
                            wire:click="quickAdd({{ $student->id }}, 'homework_missing')">
                        </flux:button>
                    </flux:button.group>
                </div>
            @endforeach
        </div>
    </flux:card>

    <!-- 自定义加减分面板 -->
    @if($selectedStudentId)
        <flux:card class="bc-card mt-6">
            <flux:heading size="md" class="mb-4">
                🎯 自定义加减分 - {{ Student::find($selectedStudentId)?->name }}
            </flux:heading>
            
            <div class="grid md:grid-cols-2 gap-4">
                <flux:select wire:model.live="scoreType" label="操作类型">
                    <flux:select.option value="add">加分</flux:select.option>
                    <flux:select.option value="subtract">减分</flux:select.option>
                </flux:select>
                
                <flux:input wire:model.live="amount" label="积分数量" type="number" min="1" max="100" />
                
                <flux:select wire:model.live="ruleCode" label="积分规则">
                    @foreach($rules as $code => $rule)
                        <flux:select.option value="{{ $code }}">
                            {{ $rule['name'] }} {{ $rule['type'] === 'add' ? '+' : '-' }}{{ $rule['amount'] }}
                        </flux:select.option>
                    @endforeach
                </flux:select>
                
                <flux:input wire:model.live="comment" label="备注" placeholder="可选备注" />
            </div>
            
            <div class="mt-4 flex gap-3">
                <flux:button wire:click="customScore" variant="primary" icon="check">
                    确认操作
                </flux:button>
                <flux:button wire:click="$set('selectedStudentId', null)" variant="ghost">
                    取消
                </flux:button>
            </div>
        </flux:card>
    @endif

    <!-- 批量加分模态框 -->
    <flux:modal wire:model.live="showBatchModal" name="batch-modal" size="md">
        <flux:heading size="lg" class="mb-4">👥 批量加分</flux:heading>
        
        <div class="space-y-4">
            <flux:input wire:model.live="batchAmount" label="每个学生加分数量" type="number" />
            <flux:select wire:model.live="batchRuleCode" label="积分规则">
                @foreach($rules as $code => $rule)
                    @if($rule['type'] === 'add')
                        <flux:select.option value="{{ $code }}">{{ $rule['name'] }}</flux:select.option>
                    @endif
                @endforeach
            </flux:select>
            <flux:input wire:model.live="batchComment" label="备注" />
            
            <flux:text class="opacity-60">
                已选择 {{ count($selectedStudents) }} 名学生
            </flux:text>
        </div>
        
        <div class="mt-6 flex gap-3">
            <flux:button wire:click="batchAdd" variant="primary">
                确认批量加分
            </flux:button>
            <flux:button wire:click="$set('showBatchModal', false)" variant="ghost">
                取消
            </flux:button>
        </div>
    </flux:modal>
</div>
