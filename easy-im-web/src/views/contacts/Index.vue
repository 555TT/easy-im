<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import FriendList from './components/FriendList.vue'
import FriendDetail from './components/FriendDetail.vue'
import FriendRequestList from './components/FriendRequestList.vue'
import AddFriendDialog from './components/AddFriendDialog.vue'
import { useContactStore } from '@/stores/contact'

const route = useRoute()
const contact = useContactStore()
const addOpen = ref(false)

const friendId = computed(() => (route.params.friendId as string) || '')
const showRequests = computed(() => friendId.value === 'requests')

onMounted(() => { contact.fetchAll() })
</script>

<template>
  <FriendList v-model:add-open="addOpen" />
  <FriendRequestList v-if="showRequests" />
  <FriendDetail v-else :friend-id="friendId" />
  <AddFriendDialog v-model:open="addOpen" />
</template>
