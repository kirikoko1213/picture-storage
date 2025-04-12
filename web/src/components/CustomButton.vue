<template>
  <button
    class="custom-button"
    :class="[theme, { 'is-disabled': disabled }]"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot name="icon">
      <span v-if="icon" class="icon">
        <component :is="icon" />
      </span>
    </slot>
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  theme: {
    type: String,
    default: 'candy',
    validator: (value: string) => ['candy', 'macaron', 'cute'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  },
  icon: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['click'])

const handleClick = (event: MouseEvent) => {
  if (!props.disabled) {
    emit('click', event)
  }
}
</script>

<style scoped>
.custom-button {
  height: 38px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 14px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  outline: none;
}

/* 糖果主题 */
.candy {
  background-color: #ff6b9c;
  color: white;
  box-shadow: 0 4px 15px rgba(255, 107, 156, 0.3);
}

.candy:hover {
  background-color: #ff4d8c;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 107, 156, 0.4);
}

.candy:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(255, 107, 156, 0.2);
}

/* 马卡龙主题 */
.macaron {
  background-color: #a8e6cf;
  color: #333;
  box-shadow: 0 4px 15px rgba(168, 230, 207, 0.3);
}

.macaron:hover {
  background-color: #8cd9b9;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(168, 230, 207, 0.4);
}

.macaron:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(168, 230, 207, 0.2);
}

/* 可爱主题 */
.cute {
  background-color: #ffd6e7;
  color: #ff6b9c;
  box-shadow: 0 4px 15px rgba(255, 214, 231, 0.3);
}

.cute:hover {
  background-color: #ffc2dc;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 214, 231, 0.4);
}

.cute:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(255, 214, 231, 0.2);
}

/* 禁用状态 */
.is-disabled {
  opacity: 0.6;
  cursor: not-allowed;
  pointer-events: none;
}

.icon {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style> 