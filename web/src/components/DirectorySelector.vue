<template>
  <div>
    <CustomButton theme="candy" @click="showDirectoryDialog = true">
      <el-icon><Folder /></el-icon>
      选择目录
    </CustomButton>

    <el-dialog
      v-model="showDirectoryDialog"
      title="选择目录"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="directory-container">
        <div
          v-for="item in directoryData"
          :key="item"
          class="directory-item"
          :class="{ active: selectedDirectory === item }"
          @click="handleDirectorySelect(item)"
        >
          <el-icon><Folder /></el-icon>
          <span>{{ item }}</span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { Folder } from "@element-plus/icons-vue"
import { apiGetDirectoryList } from "@/api-service/image-manage"
import CustomButton from "./CustomButton.vue";

const selectedDirectory = ref("")

const emit = defineEmits<{
  (e: "changeDirectory", value: string): void
}>()

const showDirectoryDialog = ref(false)
const directoryData = ref<string[]>([])

const handleDirectorySelect = (directory: string) => {
  selectedDirectory.value = directory
  showDirectoryDialog.value = false
}

onMounted(() => {
  apiGetDirectoryList().then((res) => {
    directoryData.value = res.data || []
    if (res.data && res.data.length > 0) {
      selectedDirectory.value = res.data[0]
      emit("changeDirectory", res.data[0])
    }
  })
})
</script>

<style scoped>
.directory-container {
  padding: 8px 0;
}

.directory-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 10px;
  margin-bottom: 5px;
}

.directory-item:hover {
  background-color: #f0f2f5;
}

.directory-item.active {
  background-color: #ffd6e7;
  color: #ff6b9c;
}

:deep(.el-dialog) {
  border-radius: 20px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background-color: #ffd6e7;
  margin: 0;
  padding: 15px 20px;
}

:deep(.el-dialog__title) {
  color: #ff6b9c;
  font-weight: 600;
}

:deep(.el-dialog__body) {
  padding: 20px;
}
</style>
