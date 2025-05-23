<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import type { ICreateWalletRequest } from '~/types'
import { toast } from 'vue-sonner'
import { getErrorMessage } from '~/lib/utils'

definePageMeta({
  layout: 'settings'
})

const router = useRouter()
const { 
  createWallet: mutateCreateWallet,
  isCreating, 
  error: mutationError 
} = useWalletMutations()

const formData = reactive<ICreateWalletRequest>({
  name: '',
  chainType: '',
  tags: {}
})
const tagsList = ref([{ key: '', value: '' }])

const addTag = () => {
  tagsList.value.push({ key: '', value: '' })
}

const removeTag = (index: number) => {
  tagsList.value.splice(index, 1)
}

watch(mutationError, (newError) => {
  if (newError) {
    toast.error(getErrorMessage(newError, 'An unknown error occurred while creating the wallet.'))
  }
})

const handleSubmit = async () => {
  mutationError.value = null

  // Convert tagsList to the Record<string, string> format
  const tags: Record<string, string> = tagsList.value
    .filter(tag => tag.key.trim() !== '' && tag.value.trim() !== '') // Filter out empty tags
    .reduce((acc, tag) => {
      acc[tag.key.trim()] = tag.value.trim()
      return acc
    }, {} as Record<string, string>)

  const payload: ICreateWalletRequest = {
    name: formData.name.trim(),
    chainType: formData.chainType.trim(),
    tags: Object.keys(tags).length > 0 ? tags : undefined // Only send tags if not empty
  }

  // Basic validation
  if (!payload.name || !payload.chainType) {
    toast.error('Wallet Name and Chain Type are required.')
    return
  }

  const newWallet = await mutateCreateWallet(payload)

  if (newWallet) {
    toast.success('Wallet created successfully!')
    router.push('/settings/wallets')
  }
}
</script>

<template>
  <div class="flex justify-center">
    <Card class="w-full max-w-2xl">
      <CardHeader>
        <CardTitle>Create New Wallet</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="space-y-6" @submit.prevent="handleSubmit">
          <div class="space-y-2">
            <Label for="name">Wallet Name</Label>
            <Input id="name" v-model="formData.name" required placeholder="My Ethereum Wallet" />
          </div>

          <div class="space-y-2">
            <Label for="chainType">Chain Type</Label>
            <ChainSelect 
              id="chainType" 
              v-model="formData.chainType" 
              required 
            />
          </div>

          <div class="space-y-4">
            <Label>Tags (Optional)</Label>
            <div v-for="(tag, index) in tagsList" :key="index" class="flex items-center gap-2">
              <Input v-model="tag.key" placeholder="Key" class="flex-1" />
              <Input v-model="tag.value" placeholder="Value" class="flex-1" />
              <Button type="button" variant="outline" size="icon" :disabled="tagsList.length <= 1" @click="removeTag(index)">
                <Icon name="lucide:trash-2" class="h-4 w-4" />
              </Button>
            </div>
            <Button type="button" variant="outline" size="sm" @click="addTag">
              <Icon name="lucide:plus" class="h-4 w-4 mr-1" />
              Add Tag
            </Button>
          </div>
        </form>
      </CardContent>
      <CardFooter class="flex justify-end gap-2">
         <NuxtLink to="/settings/wallets">
            <Button variant="outline">Cancel</Button>
          </NuxtLink>
        <Button type="submit" :disabled="isCreating" @click="handleSubmit">
          <Icon v-if="isCreating" name="svg-spinners:3-dots-fade" class="w-4 h-4 mr-2" />
          {{ isCreating ? 'Creating...' : 'Create Wallet' }}
        </Button>
      </CardFooter>
    </Card>
  </div>
</template>
