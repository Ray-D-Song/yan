<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue'
import { ChevronRight, Plus } from "lucide-vue-next"
import { notesApi, type Note } from '@/api/notes'
import { arrayToTree, type TreeNode } from '@/lib/utils'

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from '@/components/ui/sidebar'
import { Input } from '@/components/ui/input'
import { toast } from 'vue-sonner'
import { useRouter } from 'vue-router'

const newNote = ref({
  visible: false,
  title: ''
})

const inputRef = ref<InstanceType<typeof Input> | null>(null)
const isSubmitting = ref(false)

async function handleNewNote() {
  newNote.value.title = ''
  newNote.value.visible = true
  await nextTick()
  inputRef.value?.$el?.focus()
}

async function submitNewNote() {
  if (isSubmitting.value || !newNote.value.title.trim()) {
    newNote.value.visible = false
    newNote.value.title = ''
    return
  }

  isSubmitting.value = true
  try {
    const createdNote = await notesApi.create({
      title: newNote.value.title.trim(),
      parentId: null,
    })

    // Add to workspaces array
    workspaces.value.push({
      ...createdNote,
      children: []
    })

    toast.success('Note created successfully')
    newNote.value.visible = false
    newNote.value.title = ''
  } catch (error) {
    console.error(error)
    toast.error('Failed to create note')
  } finally {
    isSubmitting.value = false
  }
}

function handleInputKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter') {
    submitNewNote()
  } else if (event.key === 'Escape') {
    newNote.value.visible = false
    newNote.value.title = ''
  }
}

const workspaces = ref<TreeNode<Note>[]>([])

onMounted(async () => {
  try {
    const notes = await notesApi.listAll()
    workspaces.value = arrayToTree(notes, {
      sortBy: (a, b) => a.position - b.position
    })
  } catch (error) {
    console.error(error)
    toast.error('Failed to list the notes')
  }
})

const router = useRouter()
function handleClickItem(item: Note) {
  router.push(`/edit/${item.id}`)
}
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel>Workspaces</SidebarGroupLabel>
    <SidebarGroupContent>
      <SidebarMenu>
        <Collapsible v-for="workspace in workspaces" :key="workspace.id">
          <SidebarMenuItem>
            <SidebarMenuButton as-child @click="handleClickItem(workspace)">
              <a href="#">
                <span v-if="workspace.icon">{{ workspace.icon }}</span>
                <span>{{ workspace.title }}</span>
              </a>
            </SidebarMenuButton>
            <CollapsibleTrigger v-if="workspace.children && workspace.children.length > 0" as-child>
              <SidebarMenuAction
                class="left-2 bg-sidebar-accent text-sidebar-accent-foreground data-[state=open]:rotate-90"
                show-on-hover
              >
                <ChevronRight />
              </SidebarMenuAction>
            </CollapsibleTrigger>
            <SidebarMenuAction show-on-hover>
              <Plus />
            </SidebarMenuAction>
            <CollapsibleContent v-if="workspace.children && workspace.children.length > 0">
              <SidebarMenuSub>
                <SidebarMenuSubItem v-for="page in workspace.children" :key="page.id">
                  <SidebarMenuSubButton as-child>
                    <a href="#">
                      <span v-if="page.icon">{{ page.icon }}</span>
                      <span>{{ page.title }}</span>
                    </a>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>

        <SidebarMenuItem v-if="newNote.visible">
          <SidebarMenuButton class="text-sidebar-foreground/70">
            <Input
              ref="inputRef"
              v-model="newNote.title"
              placeholder="Enter note title..."
              :disabled="isSubmitting"
              @blur="submitNewNote"
              @keydown="handleInputKeydown"
            />
          </SidebarMenuButton>
        </SidebarMenuItem>

        <SidebarMenuItem>
          <SidebarMenuButton class="text-sidebar-foreground/70" @click="handleNewNote">
            <Plus />
            <span>New</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
