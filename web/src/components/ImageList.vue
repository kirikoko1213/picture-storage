<template>
  <div class="image-list">
    <div class="toolbar">
      <div class="left">
        
        <DirectorySelector @change-directory="handleChangeDirectory" />
        <CustomButton theme="candy" @click="handleSearch">
            <template #icon>
                <SearchIcon color="#fff" :size="14"/>
            </template>
            <template #default>
                搜索
            </template>
        </CustomButton>
        <el-upload
          multiple
          :show-file-list="false"
          accept="image/*"
          action="/api/upload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :data="uploadData"
          :before-upload="handleBeforeUpload"
        >
          <CustomButton theme="candy">
            <template #icon>
                <UploadIcon color="#fff" :size="14"/>
            </template>
            <template #default>
                上传
            </template>
          </CustomButton>
        </el-upload>
        <CustomButton theme="candy" @click="handleBatchDelete">
            <template #icon>
                <DeleteIcon color="#fff" :size="15"/>
            </template>
            <template #default>
                批量删除
            </template>
        </CustomButton>
      </div>
      <div class="right">
        <div>
          <CustomSelect
            theme="candy"
            @change="(value) => (selectedTags = value)"
            multiple
            filterable
            allow-create
            placeholder="请选择或输入标签"
            style="width: 100%"
            :options="tagOptions.map((tag) => ({ label: tag, value: tag }))"
          >
          </CustomSelect>
        </div>
        <div
          class="icon-container"
          :class="{ active: viewMode === 'list' }"
          @click="viewMode = 'list'"
        >
          <ListIcon color="#ff6b9c" />
        </div>
        <div
          class="icon-container"
          :class="{ active: viewMode === 'grid' }"
          @click="viewMode = 'grid'"
        >
          <GridIcon color="#ff6b9c" />
        </div>
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
            :preview-src-list="imageURLList"
            :initial-index="row.index"
            fit="cover"
            preview-teleported
          />
        </template>
      </el-table-column>
      <el-table-column prop="size" label="大小">
        <template #default="{ row }">
          {{ formatFileSize(row.size) }}
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="上传时间">
        <template #default="{ row }">
          {{ dayjs(row.createdAt).format("YYYY-MM-DD HH:mm:ss") }}
        </template>
      </el-table-column>
      <el-table-column prop="tags" label="标签">
        <template #default="{ row }">
          <div class="tags">
            <el-tag v-for="tag in row.tags" :key="tag" type="success">{{ tag }}</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 网格视图 -->
    <div v-else class="grid-view">
      <div v-for="image in imageList" :key="image.id" class="grid-item">
        <el-image
          :src="image.url"
          :preview-src-list="imageURLList"
          :initial-index="image.index"
          fit="cover"
          class="grid-image"
        />
        <div class="image-info">
          <div class="image-name">{{ image.name }}</div>
          <div class="image-size">{{ formatFileSize(image.size) }}</div>
        </div>
      </div>
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
import { ref, onMounted, computed } from "vue"
import { dayjs, ElMessage } from "element-plus"
import { apiGetImageList, apiDeleteImages, apiGetTags } from "@/api-service/image-manage"
import type { ImageItem } from "@/api-service/image-manage"
import { formatFileSize } from "@/helper/file"
import GridIcon from "@/icons/GridIcon.vue"
import ListIcon from "@/icons/ListIcon.vue"
import DirectorySelector from "./DirectorySelector.vue"
import CustomButton from "./CustomButton.vue"
import CustomSelect from "./CustomSelect.vue"
import SearchIcon from "@/icons/SearchIcon.vue"
import UploadIcon from "@/icons/UploadIcon.vue"
import DeleteIcon from "@/icons/DeleteIcon.vue"

const selectedDirectory = ref("")

const viewMode = ref<"list" | "grid">("grid")
const imageList = ref<(ImageItem & { index: number })[]>([])
const selectedImages = ref<number[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedTags = ref<string[]>([])
const tagOptions = ref<string[]>([])
const uploadData = ref({
  tags: "",
  directory: "",
})

const fetchImageList = async () => {
  try {
    const response = await apiGetImageList(
      selectedDirectory.value || "",
      selectedTags.value,
      currentPage.value,
      pageSize.value,
    )
    if (response.data) {
      imageList.value = response.data.list?.map((item, index) => ({ ...item, index })) || []
      total.value = response.data.total
    }
  } catch (error) {
    ElMessage.error("获取图片列表失败")
  }
}

const handleChangeDirectory = (directory: string) => {
  selectedDirectory.value = directory
  fetchImageList()
}

const imageURLList = computed(() => {
  return imageList.value.map((item) => item.url)
})

const handleSearch = () => {
  fetchImageList()
}

const handleBeforeUpload = (file: File) => {
  uploadData.value.tags = selectedTags.value?.join(",") || ""
  uploadData.value.directory = selectedDirectory.value
  return true
}

const handleUploadSuccess = (response: any) => {
  ElMessage.success("上传成功")
  selectedTags.value = []
  fetchImageList()
}

const handleUploadError = () => {
  ElMessage.error("上传失败")
}

const handleBatchDelete = async () => {
  if (selectedImages.value.length === 0) {
    ElMessage.warning("请选择要删除的图片")
    return
  }
  try {
    await apiDeleteImages(selectedImages.value)
    ElMessage.success("删除成功")
    fetchImageList()
  } catch (error) {
    ElMessage.error("删除失败")
  }
}

const handleDelete = async (image: ImageItem) => {
  try {
    await apiDeleteImages([image.id])
    ElMessage.success("删除成功")
    fetchImageList()
  } catch (error) {
    ElMessage.error("删除失败")
  }
}

const handleSelectionChange = (selection: ImageItem[]) => {
  selectedImages.value = selection.map((item) => item.id)
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchImageList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchImageList()
}

const fetchTagOptions = async () => {
  const response = await apiGetTags()
  if (response.data) {
    tagOptions.value = response.data || []
  }
}

onMounted(() => {
  fetchImageList()
  fetchTagOptions()
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
  .left {
    display: flex;
    gap: 10px;
  }
  .right {
    display: flex;
    gap: 10px;
    .icon-container {
      cursor: pointer;
      display: flex;
      align-items: center;
      padding: 4px;
      border-radius: 10px;
      transition: all 0.3s ease;
      svg {
        color: #bac6d2;
      }
      &:hover {
        background-color: #f0f0f0;
      }
    }
    .active {
      background-color: #ffd6e7;
      svg {
        color: #ff6b9c;
      }
    }
  }
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
  border-radius: 10px;
}

.image-info {
  padding: 10px;
}

.image-name {
  font-size: 14px;
  margin-bottom: 5px;
  color: #333;
  font-weight: 500;
}

.image-size {
  font-size: 12px;
  color: #999;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  .el-pagination {
    --el-pagination-button-bg-color: #fff;
    --el-pagination-hover-color: #ff6b9c;
    .el-pager li {
      border-radius: 10px;
      &.is-active {
        background-color: #ff6b9c;
      }
    }
  }
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  .el-tag {
    border-radius: 10px;
    padding: 0 10px;
    height: 24px;
    line-height: 24px;
    background-color: #ffd6e7;
    border-color: #ffd6e7;
    color: #ff6b9c;
  }
}

.el-table {
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  .el-table__header {
    th {
      background-color: #ffd6e7;
      color: #ff6b9c;
      font-weight: 600;
    }
  }
  .el-table__row {
    &:hover {
      background-color: #fff5f9;
    }
  }
  .el-button {
    border-radius: 10px;
    padding: 6px 12px;
  }
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
