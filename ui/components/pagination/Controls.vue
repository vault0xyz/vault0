<script setup lang="ts">
import { computed } from 'vue';

interface PageControlsProps {
  nextToken?: string;
  currentToken?: string;
}

const props = defineProps<PageControlsProps>();

// Define the events emitted by this component
const emit = defineEmits<{ 
  (e: 'previous' | 'next'): void 
}>();

const isFirstPage = computed(() => !props.currentToken);

// Determine if there's a next page based on nextToken
const hasNextPage = computed(() => !!props.nextToken && props.nextToken.trim() !== '');
</script>

<template>
  <div class="flex items-center space-x-2">
    <Button 
      variant="outline" 
      size="icon" 
      :disabled="isFirstPage" 
      @click="!isFirstPage && emit('previous')"
    >
      <span class="sr-only">First page</span>
      <Icon name="lucide:chevrons-left" class="h-4 w-4" />
    </Button>
    <Button 
      variant="outline" 
      size="icon" 
      :disabled="!hasNextPage" 
      @click="hasNextPage && emit('next')"
    >
      <span class="sr-only">Next page</span>
      <Icon name="lucide:chevron-right" class="h-4 w-4" />
    </Button>
  </div>
</template> 