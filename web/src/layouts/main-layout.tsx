import { defineComponent } from 'vue'
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
import { RouterView } from 'vue-router'

export default defineComponent({
  name: 'MainLayout',
  setup() {
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
                      Yan - Note Taking App
                    </BreadcrumbPage>
                  </BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
            </div>
          </header>
          <div class="flex flex-1 flex-col">
            <RouterView />
          </div>
        </SidebarInset>
        <SidebarRight />
      </SidebarProvider>
    )
  },
})
