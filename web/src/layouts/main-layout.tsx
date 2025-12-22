import { defineComponent, ref, watch, onUnmounted } from 'vue'
import SidebarLeft from '@/components/SidebarLeft.vue'
import SidebarRight from '@/components/SidebarRight.vue'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from '@/components/ui/breadcrumb'
import { Separator } from '@/components/ui/separator'
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { MilkdownProvider } from '@milkdown/vue'
import MdEditor from '@/components/app/md-editor'
import { useRoute } from 'vue-router'
import { notesApi, type Note } from '@/api/notes'

export default defineComponent({
  name: 'MainLayout',
  setup() {
    const route = useRoute()
    const currentNote = ref<Note | null>(null)
    const loading = ref(false)
    const editingContent = ref('')
    const noteId = ref<number | null>(null)
    const needSync = ref(false)
    const syncing = ref(false)

    watch(() => route.query, async (query) => {
      const id = query.id
      if (id && typeof id === 'string') {
        try {
          loading.value = true
          const note = await notesApi.getById(Number(id))
          currentNote.value = note
          noteId.value = note.id
          editingContent.value = note.content
        } catch (error) {
          console.error('Failed to fetch note:', error)
          currentNote.value = null
          noteId.value = null
          editingContent.value = ''
        } finally {
          loading.value = false
        }
      } else {
        currentNote.value = null
        noteId.value = null
        editingContent.value = ''
      }
    }, { immediate: true })

    const syncContent = async () => {
      if (!needSync.value || syncing.value || !noteId.value || !currentNote.value) {
        return
      }

      try {
        syncing.value = true
        await notesApi.update(noteId.value, {
          title: currentNote.value.title,
          content: editingContent.value,
        })
        needSync.value = false
      } catch (error) {
        console.error('Failed to sync note:', error)
      } finally {
        syncing.value = false
      }
    }

    // Check for sync every 3 seconds
    const syncInterval = setInterval(() => {
      syncContent()
    }, 3000)

    onUnmounted(() => {
      clearInterval(syncInterval)
    })

    const handleUpdateContent = (value: string) => {
      editingContent.value = value
      needSync.value = true
    }

    return () => (
      <SidebarProvider>
        <SidebarLeft />
        <SidebarInset>
          <header class="sticky top-0 flex h-14 shrink-0 items-center gap-2 bg-background">
            <div class="flex flex-1 items-center gap-2 px-3">
              <SidebarTrigger />
              <Separator
                orientation="vertical"
                class="mr-2 data-[orientation=vertical]:h-4"
              />
              <Breadcrumb>
                <BreadcrumbList>
                  <BreadcrumbItem>
                    <BreadcrumbPage class="line-clamp-1">
                      Project Management & Task Tracking
                    </BreadcrumbPage>
                  </BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
            </div>
          </header>
          <div class="flex flex-1 flex-col gap-4 p-4">
          <MilkdownProvider>
            <MdEditor
              id={noteId.value}
              value={editingContent.value}
              onUpdateValue={handleUpdateContent}
            />
          </MilkdownProvider>
          </div>
        </SidebarInset>
        <SidebarRight />
      </SidebarProvider>
    )
  },
})
