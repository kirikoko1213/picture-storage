<template>
  <div class="common-layout">
    <el-container>
      <el-main>
        <div class="tag-manage">
          <div class="header">
            <h2>标签管理</h2>
            <el-button type="primary" @click="showAddDialog = true">
              添加标签
            </el-button>
          </div>
          
          <el-table 
            :data="tagList" 
            v-loading="loading"
            style="width: 100%"
          >
            <el-table-column prop="name" label="标签名" />
            <el-table-column prop="count" label="图片数量" />
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="editTag(scope.row)"
                >
                  编辑
                </el-button>
                <el-button 
                  type="danger" 
                  size="small" 
                  @click="deleteTag(scope.row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-main>
    </el-container>

    <!-- 添加标签对话框 -->
    <el-dialog v-model="showAddDialog" title="添加标签" width="400px">
      <el-form :model="addForm" label-width="80px">
        <el-form-item label="标签名">
          <el-input v-model="addForm.name" placeholder="请输入标签名" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取消</el-button>
          <el-button type="primary" @click="handleAddTag">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑标签对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑标签" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="标签名">
          <el-input v-model="editForm.name" placeholder="请输入标签名" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="handleEditTag">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { apiGetTagDetails, apiCreateTag, apiDeleteTag, apiUpdateTag } from '@/api-service/image-manage'
import type { TagItem } from '@/api-service/image-manage'

const tagList = ref<TagItem[]>([])
const loading = ref(false)
const showAddDialog = ref(false)
const showEditDialog = ref(false)

const addForm = ref({
  name: ''
})

const editForm = ref({
  name: '',
  oldName: ''
})

// 获取标签列表
const fetchTags = async () => {
  loading.value = true
  try {
    const response = await apiGetTagDetails()
    if (response.status === 'success') {
      tagList.value = response.data?.list || []
    } else {
      ElMessage.error(response.msg || '获取标签列表失败')
    }
  } catch (error) {
    ElMessage.error('获取标签列表失败')
  } finally {
    loading.value = false
  }
}

// 添加标签
const handleAddTag = async () => {
  if (!addForm.value.name.trim()) {
    ElMessage.warning('请输入标签名')
    return
  }
  
  try {
    const response = await apiCreateTag(addForm.value.name.trim())
    if (response.status === 'success') {
      ElMessage.success('添加标签成功')
      showAddDialog.value = false
      addForm.value.name = ''
      fetchTags()
    } else {
      ElMessage.error(response.msg || '添加标签失败')
    }
  } catch (error) {
    ElMessage.error('添加标签失败')
  }
}

// 编辑标签
const editTag = (tag: TagItem) => {
  editForm.value.name = tag.name
  editForm.value.oldName = tag.name
  showEditDialog.value = true
}

// 处理编辑标签
const handleEditTag = async () => {
  if (!editForm.value.name.trim()) {
    ElMessage.warning('请输入标签名')
    return
  }
  
  if (editForm.value.name === editForm.value.oldName) {
    ElMessage.warning('标签名没有变化')
    return
  }
  
  try {
    const response = await apiUpdateTag(editForm.value.oldName, editForm.value.name.trim())
    if (response.status === 'success') {
      ElMessage.success('编辑标签成功')
      showEditDialog.value = false
      fetchTags()
    } else {
      ElMessage.error(response.msg || '编辑标签失败')
    }
  } catch (error) {
    ElMessage.error('编辑标签失败')
  }
}

// 删除标签
const deleteTag = async (tag: TagItem) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除标签 "${tag.name}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await apiDeleteTag(tag.name)
    if (response.status === 'success') {
      ElMessage.success('删除标签成功')
      fetchTags()
    } else {
      ElMessage.error(response.msg || '删除标签失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除标签失败')
    }
  }
}

onMounted(() => {
  fetchTags()
})
</script>

<style scoped>
.common-layout {
  height: 100vh;
}

.tag-manage {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}

.dialog-footer {
  display: flex;
  gap: 10px;
}
</style> 