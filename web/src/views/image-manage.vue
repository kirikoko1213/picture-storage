<template>
  <div class="common-layout">
    <el-container>
      <el-main style="height: calc(100vh - 150px);">
        <ImageList />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import ImageList from '@/components/ImageList.vue'

</script>

<style scoped>
.common-layout {
  height: 100vh;
}
</style>