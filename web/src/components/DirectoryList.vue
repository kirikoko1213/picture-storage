<template>
  <div class="directory-list">
      <div class="directory-container">
        <div
          v-for="item in directoryData"
          :key="item"
          class="directory-item"
          :class="{ active: activeIndex === item }"
          @click="handleSelect(item)"
        >
          <el-icon><Folder /></el-icon>
          <span>{{ item }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Folder } from '@element-plus/icons-vue'
import { apiGetDirectoryList } from '@/api-service/image-manage'

const directoryData = ref<string[]>([])

const emit = defineEmits<{
  (e: 'select', name: string): void
}>()

const activeIndex = ref('')

const handleSelect = (index: string) => {
  activeIndex.value = index
  emit('select', index)
}

onMounted(() => {
  apiGetDirectoryList().then((res) => {
    directoryData.value = res.data || []
    if (directoryData.value.length > 0) {
      emit('select', directoryData.value[0])
    }
  })
})
</script>

<style scoped>
.directory-list {
  height: 100%;
  background-color: #f5f7fa;
}

.directory-container {
  padding: 8px 0;
}

.directory-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.directory-item:hover {
  background-color: #f0f2f5;
}

.directory-item.active {
  background-color: #ecf5ff;
  color: #409eff;
}
</style>