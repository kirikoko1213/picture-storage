<template>
  <div class="image-list">
    <div class="toolbar">
      <div class="left">
        <el-upload
          class="upload-demo"
          action="/api/upload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :data="uploadData"
          :before-upload="handleBeforeUpload"
        >
          <el-button type="primary">上传图片</el-button>
          <template #tip>
            <div class="el-upload__tip">
              <el-select
                v-model="selectedTags"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="请选择或输入标签"
                style="width: 100%"
              >
                <el-option
                  v-for="tag in tagOptions"
                  :key="tag"
                  :label="tag"
                  :value="tag"
                />
              </el-select>
            </div>
          </template>
        </el-upload>
        <el-button type="danger" @click="handleBatchDelete">批量删除</el-button>
      </div>
      <div class="right">
        <el-radio-group v-model="viewMode" size="small">
          <el-radio-button label="list">列表视图</el-radio-button>
          <el-radio-button label="grid">网格视图</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <!-- 列表视图 -->
    <el-table
      v-if="viewMode === 'list'"
      :data="imageList"
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="imageName" label="文件名" />
      <el-table-column prop="url" label="预览">
        <template #default="{ row }">
          <el-image
            style="width: 100px; height: 100px"
            :src="row.url"
            :preview-src-list="[row.url]"
            fit="cover"
          />
        </template>
      </el-table-column>
      <el-table-column prop="size" label="大小">
        <template #default="{ row }">
          {{ formatFileSize(row.size) }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="上传时间" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 网格视图 -->
    <div v-else class="grid-view">
      <el-checkbox-group v-model="selectedImages">
        <div v-for="image in imageList" :key="image.id" class="grid-item">
          <el-checkbox :label="image.id">
            <el-image
              :src="image.url"
              :preview-src-list="[image.url]"
              fit="cover"
              class="grid-image"
            />
            <div class="image-info">
              <div class="image-name">{{ image.name }}</div>
              <div class="image-size">{{ formatFileSize(image.size) }}</div>
            </div>
          </el-checkbox>
        </div>
      </el-checkbox-group>
    </div>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { apiGetImageList, apiDeleteImages } from '@/api-service/image-manage'
import type { ImageItem } from '@/api-service/image-manage'

const props = defineProps({
  directory: {
    type: String,
    required: true
  }
})

const viewMode = ref<'list' | 'grid'>('list')
const imageList = ref<ImageItem[]>([])
const selectedImages = ref<number[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedTags = ref<string[]>([])
const tagOptions = ref<string[]>([])
const uploadData = ref({
  tags: '',
  directory: ''
})

const fetchImageList = async () => {
  try {
    const response = await apiGetImageList(props.directory, currentPage.value, pageSize.value)
    if (response.data) {
      imageList.value = response.data.list
      total.value = response.data.total
    }
  } catch (error) {
    ElMessage.error('获取图片列表失败')
  }
}

const handleBeforeUpload = (file: File) => {
  uploadData.value.tags = selectedTags.value?.join(',') || ''
  uploadData.value.directory = props.directory
  return true
}

const handleUploadSuccess = (response: any) => {
  ElMessage.success('上传成功')
  selectedTags.value = []
  fetchImageList()
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

const handleBatchDelete = async () => {
  if (selectedImages.value.length === 0) {
    ElMessage.warning('请选择要删除的图片')
    return
  }
  try {
    await apiDeleteImages(selectedImages.value)
    ElMessage.success('删除成功')
    fetchImageList()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleDelete = async (image: ImageItem) => {
  try {
    await apiDeleteImages([image.id])
    ElMessage.success('删除成功')
    fetchImageList()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSelectionChange = (selection: ImageItem[]) => {
  selectedImages.value = selection.map(item => item.id)
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchImageList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchImageList()
}

const formatFileSize = (size: number) => {
  if (size < 1024) {
    return size + ' B'
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + ' KB'
  } else if (size < 1024 * 1024 * 1024) {
    return (size / (1024 * 1024)).toFixed(2) + ' MB'
  } else {
    return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
  }
}

onMounted(() => {
  fetchImageList()
})
</script>

<style scoped>
.image-list {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.left {
  display: flex;
  gap: 10px;
}

.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.grid-item {
  position: relative;
}

.grid-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.image-info {
  padding: 10px;
}

.image-name {
  font-size: 14px;
  margin-bottom: 5px;
}

.image-size {
  font-size: 12px;
  color: #999;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style> 