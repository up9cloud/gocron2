<script setup>
import breadcrumb from '../components/breadcrumb.vue'
</script>
<template>
  <el-container>
    <el-header>
      <slot name="header">
        <breadcrumb></breadcrumb>
      </slot>
    </el-header>
    <slot></slot>
    <el-footer>
      <slot name="footer"></slot>
    </el-footer>
  </el-container>
</template>
