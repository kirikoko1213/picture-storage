<template>
  <div class="app">
    <div class="header">
      <div class="tab-container">
        <div class="custom-tabs">
          <div 
            class="tab-item" 
            :class="{ active: activeTab === 'images' }"
            @click="switchTab('images')"
          >
            <i class="tab-icon">🖼️</i>
            <span>图片管理</span>
          </div>
          <div 
            class="tab-item" 
            :class="{ active: activeTab === 'tags' }"
            @click="switchTab('tags')"
          >
            <i class="tab-icon">🏷️</i>
            <span>标签管理</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="content">
      <!-- 使用 v-if 确保组件在切换时销毁重新创建 -->
      <ImageManage v-if="activeTab === 'images'" :key="'images-' + tabKey" />
      <TagManage v-if="activeTab === 'tags'" :key="'tags-' + tabKey" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ImageManage from './views/image-manage.vue'
import TagManage from './views/tag-manage.vue'

const activeTab = ref('images')
const tabKey = ref(0)

const switchTab = (tab: string) => {
  activeTab.value = tab
  // 增加 key 值，强制组件重新创建
  tabKey.value++
}
</script>

<style scoped>
.app {
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.header {
  border-radius: 10px 10px 0 0;
  background: linear-gradient(135deg, #e96692 0%, #ffdae2 100%);
  padding: 20px 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.title {
  color: white;
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 20px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  letter-spacing: 1px;
}

.tab-container {
  display: flex;
  justify-content: flex-start;
}

.custom-tabs {
  display: flex;
  gap: 8px;
  background: rgba(255, 255, 255, 0.1);
  padding: 6px;
  border-radius: 12px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
  font-size: 15px;
  position: relative;
  overflow: hidden;
}

.tab-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(45deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transform: translateX(-100%);
  transition: transform 0.6s;
}

.tab-item:hover::before {
  transform: translateX(100%);
}

.tab-item:hover {
  color: white;
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-2px);
}

.tab-item.active {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.4);
  transform: translateY(-1px);
}

.tab-item.active:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 107, 107, 0.5);
}

.tab-icon {
  font-size: 18px;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.content {
  flex: 1;
  overflow: hidden;
  margin: 20px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.8);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    padding: 15px 20px;
  }
  
  .title {
    font-size: 24px;
    margin-bottom: 15px;
  }
  
  .tab-item {
    padding: 10px 16px;
    font-size: 14px;
  }
  
  .tab-icon {
    font-size: 16px;
  }
  
  .content {
    margin: 15px;
  }
}

@media (max-width: 480px) {
  .custom-tabs {
    flex-direction: column;
    gap: 4px;
  }
  
  .tab-item {
    justify-content: center;
    padding: 12px;
  }
  
  .title {
    font-size: 20px;
    text-align: center;
  }
}
</style>
