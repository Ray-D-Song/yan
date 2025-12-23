import { defineComponent, ref, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { MilkdownProvider } from '@milkdown/vue'
import MdEditor from '@/components/app/md-editor'
import { notesApi, type Note } from '@/api/notes'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from '@/components/ui/breadcrumb'

export default defineComponent({
  name: 'EditPage',
  setup() {
    const route = useRoute()
    const currentNote = ref<Note | null>(null)
    const loading = ref(false)
    const editingContent = ref('')
    const noteId = ref<number | null>(null)
    const needSync = ref(false)
    const syncing = ref(false)

    // Watch route params to load note
    watch(() => route.params.id, async (id) => {
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

    // Auto-sync content
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
      <div class="flex flex-1 flex-col">
        <header class="flex h-14 shrink-0 items-center gap-2 border-b px-4">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbPage class="line-clamp-1">
                  {currentNote.value?.title || 'Loading...'}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </header>
        <div class="flex flex-1 flex-col gap-4 p-4">
          {loading.value ? (
            <div class="flex items-center justify-center h-full">
              <span>Loading...</span>
            </div>
          ) : noteId.value ? (
            <MilkdownProvider>
              <MdEditor
                id={noteId.value}
                value={editingContent.value}
                onUpdateValue={handleUpdateContent}
              />
            </MilkdownProvider>
          ) : (
            <div class="flex items-center justify-center h-full text-muted-foreground">
              <span>Select a note to start editing</span>
            </div>
          )}
        </div>
      </div>
    )
  },
})
