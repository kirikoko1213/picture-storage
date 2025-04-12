<template>
  <div class="common-layout">
    <el-container>
      <el-container>
        <el-aside width="200px">
          <DirectoryList
            @select="handleDirectorySelect"
          />
        </el-aside>
        <el-main>
          <ImageList v-if="selectedDirectory" :directory="selectedDirectory" />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import DirectoryList from '@/components/DirectoryList.vue'
import ImageList from '@/components/ImageList.vue'

const selectedDirectory = ref<string>()

const handleDirectorySelect = (id: string) => {
  selectedDirectory.value = id
}


</script>

<style scoped>
.common-layout {
  height: 100vh;
}

.el-aside {
  background-color: #f5f7fa;
}
</style>