<template>
  <div class="custom-select">
    <el-select
      v-model="modelValue"
      :multiple="multiple"
      :filterable="filterable"
      :allow-create="allowCreate"
      :placeholder="placeholder"
      :disabled="disabled"
      :class="theme"
      @change="handleChange"
    >
      <el-option
        v-for="item in options"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </el-select>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"

interface Option {
  label: string
  value: string | number
}

const modelValue = ref<string | number | string[]>([])

const props = defineProps({
  options: {
    type: Array as () => Option[],
    default: () => [],
  },
  multiple: {
    type: Boolean,
    default: false,
  },
  filterable: {
    type: Boolean,
    default: false,
  },
  allowCreate: {
    type: Boolean,
    default: false,
  },
  placeholder: {
    type: String,
    default: "请选择",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  theme: {
    type: String,
    default: "candy",
    validator: (value: string) => ["candy", "macaron", "cute"].includes(value),
  },
})

const emit = defineEmits(["update:modelValue", "change", "search"])

const buttonTheme = computed(() => {
  return props.theme === "candy" ? "cute" : props.theme
})

const handleChange = (value: any) => {
  emit("change", value)
}

const handleSearch = () => {
  emit("search", modelValue.value)
}
</script>

<style scoped>
.custom-select {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-button {
  background-color: #fff5f9;
}

:deep(.el-select) {
  width: 100%;
}

:deep(.el-select__wrapper) {
  border-radius: 10px;
  height: 38px;
}

/* 糖果主题 */
:deep(.candy .el-select__wrapper) {
  background-color: #fff5f9;
  box-shadow: 0 0 0 1.5px #ffd6e7;
}

:deep(.candy .el-select__wrapper:hover) {
  box-shadow: 0 0 0 1.5px #ff6b9c;
}

:deep(.candy .el-select__wrapper.is-focus) {
  box-shadow: 0 0 0 1.5px #ff6b9c;
}

:deep(.candy .el-select__tags) {
  background-color: #fff5f9;
}

/* 马卡龙主题 */
:deep(.macaron .el-select__wrapper) {
  background-color: #f0f9f5;
  box-shadow: 0 0 0 1.5px #d8f3e8;
}

:deep(.macaron .el-select__wrapper:hover) {
  box-shadow: 0 0 0 1.5px #a8e6cf;
}

:deep(.macaron .el-select__wrapper.is-focus) {
  box-shadow: 0 0 0 1.5px #a8e6cf;
}

:deep(.macaron .el-select__tags) {
  background-color: #f0f9f5;
}

/* 可爱主题 */
:deep(.cute .el-select__wrapper) {
  background-color: #fff5f9;
  box-shadow: 0 0 0 1.5px #ffd6e7;
}

:deep(.cute .el-select__wrapper:hover) {
  box-shadow: 0 0 0 1.5px #ffb8d9;
}

:deep(.cute .el-select__wrapper.is-focus) {
  box-shadow: 0 0 0 1.5px #ffb8d9;
}

:deep(.cute .el-select__tags) {
  background-color: #fff5f9;
}

/* 下拉菜单样式 */
:deep(.el-select-dropdown) {
  border-radius: 10px;
  border: none;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

:deep(.el-select-dropdown__item) {
  padding: 8px 16px;
  border-radius: 6px;
  margin: 4px;
}

:deep(.el-select-dropdown__item.selected) {
  background-color: #ffd6e7;
  color: #ff6b9c;
}

:deep(.el-select-dropdown__item:hover) {
  background-color: #fff5f9;
}
</style>
