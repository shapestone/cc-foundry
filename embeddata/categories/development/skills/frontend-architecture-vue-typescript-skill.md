---
name: frontend-architecture-vue-typescript
description: Guide LLMs in implementing feature-based frontend architecture using Vue 3, TypeScript, Pinia, and Vue Router. Covers component hierarchy, Pinia store patterns, API service layers, composables, routing, cross-feature communication, and LLM-friendly coding conventions for scalable single-page applications.
---

# Frontend Architecture Skill - Vue 3 (TypeScript)

## Purpose
Guide LLMs in implementing feature-based frontend architecture using Vue 3, TypeScript, Pinia, and Vue Router.

## Core Principles

### 1. Feature-Based Organization (Mirror Backend)
- Each feature is a self-contained vertical slice
- All frontend code for a feature lives in one directory
- Mirrors backend structure for full-stack coherence
- Optimized for LLM context windows

### 2. Component Hierarchy
- Feature components → Shared components → Base components
- Components import downward, never upward
- Minimal cross-feature dependencies

### 3. State Management with Pinia
- One store per feature
- Stores are collocated with features
- Global stores live in shared/

---

## Directory Structure

```
web/src/
├── features/
│   ├── catalog/
│   │   ├── components/
│   │   │   ├── CatalogSidebar.vue
│   │   │   ├── CatalogTree.vue
│   │   │   ├── CatalogDialog.vue
│   │   │   └── EndpointCard.vue
│   │   ├── stores/
│   │   │   └── catalogStore.ts
│   │   ├── services/
│   │   │   └── catalogApi.ts
│   │   ├── types/
│   │   │   └── catalog.ts
│   │   ├── composables/
│   │   │   └── useCatalogSearch.ts
│   │   └── router/
│   │       └── routes.ts
│   │
│   ├── request/
│   │   ├── components/
│   │   │   ├── RequestForm.vue
│   │   │   ├── RequestHeader.vue
│   │   │   ├── RequestBody.vue
│   │   │   └── ResponseViewer.vue
│   │   ├── stores/
│   │   │   └── requestStore.ts
│   │   ├── services/
│   │   │   └── requestApi.ts
│   │   └── types/
│   │       └── request.ts
│   │
│   ├── invocation/
│   │   ├── components/
│   │   │   ├── InvocationHistory.vue
│   │   │   └── InvocationDetail.vue
│   │   ├── stores/
│   │   │   └── invocationStore.ts
│   │   └── services/
│   │       └── invocationApi.ts
│   │
│   ├── tab/
│   │   ├── components/
│   │   │   ├── TabBar.vue
│   │   │   └── TabContent.vue
│   │   └── stores/
│   │       └── tabStore.ts
│   │
│   ├── sandbox/
│   ├── importexport/
│   │   ├── components/
│   │   │   ├── ImportDialog.vue
│   │   │   └── ExportDialog.vue
│   │   └── services/
│   │       ├── swaggerImporter.ts
│   │       └── openApiExporter.ts
│   │
│   └── settings/
│       ├── components/
│       │   └── SettingsPanel.vue
│       └── stores/
│           └── settingsStore.ts
│
├── shared/
│   ├── components/
│   │   ├── BaseButton.vue
│   │   ├── BaseInput.vue
│   │   ├── BaseModal.vue
│   │   └── LoadingSpinner.vue
│   ├── composables/
│   │   ├── useDebounce.ts
│   │   ├── useApi.ts
│   │   └── useToast.ts
│   ├── stores/
│   │   └── featureFlagsStore.ts
│   └── types/
│       └── common.ts
│
├── router/
│   └── index.ts
│
├── App.vue
└── main.ts
```

---

## Component Architecture

### Component Responsibilities

**Feature Components** (`features/{feature}/components/`)
- Feature-specific UI logic
- Uses feature store
- Calls feature API service
- Can import shared components

**Shared Components** (`shared/components/`)
- Reusable across multiple features
- No feature-specific logic
- Generic, configurable via props
- Examples: buttons, inputs, modals

### Component Structure Template

```vue
<!-- features/catalog/components/CatalogSidebar.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useCatalogStore } from '../stores/catalogStore'
import { useCatalogSearch } from '../composables/useCatalogSearch'
import CatalogTree from './CatalogTree.vue'
import BaseButton from '@/shared/components/BaseButton.vue'

// Store
const catalogStore = useCatalogStore()
const { catalogs, isLoading } = storeToRefs(catalogStore)

// Composables
const { searchQuery, filteredCatalogs } = useCatalogSearch(catalogs)

// Methods
const handleCreateCatalog = () => {
  catalogStore.openCreateDialog()
}
</script>

<template>
  <div class="catalog-sidebar">
    <div class="sidebar-header">
      <h2>Catalogs</h2>
      <BaseButton @click="handleCreateCatalog">
        Add Catalog
      </BaseButton>
    </div>
    
    <input 
      v-model="searchQuery" 
      placeholder="Search catalogs..."
      class="search-input"
    />
    
    <div v-if="isLoading" class="loading">
      Loading...
    </div>
    
    <CatalogTree 
      v-else
      :catalogs="filteredCatalogs"
    />
  </div>
</template>

<style scoped>
.catalog-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
}

.search-input {
  margin: 0 1rem 1rem;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>
```

---

## Pinia Store Architecture

### Store Structure

**Feature Stores** (`features/{feature}/stores/`)
- One store per feature
- Manages feature-specific state
- Calls feature API service
- Can reference other stores if needed

**Shared Stores** (`shared/stores/`)
- Cross-cutting concerns
- Example: feature flags, authentication, notifications

### Store Template

```typescript
// features/catalog/stores/catalogStore.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ApiCatalog } from '../types/catalog'
import { catalogApi } from '../services/catalogApi'

export const useCatalogStore = defineStore('catalog', () => {
  // State
  const catalogs = ref<ApiCatalog[]>([])
  const selectedCatalogId = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const isCreateDialogOpen = ref(false)

  // Getters
  const selectedCatalog = computed(() => {
    if (!selectedCatalogId.value) return null
    return catalogs.value.find(c => c.id === selectedCatalogId.value)
  })

  const catalogCount = computed(() => catalogs.value.length)

  // Actions
  async function loadCatalogs() {
    isLoading.value = true
    error.value = null
    
    try {
      catalogs.value = await catalogApi.fetchAll()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to load catalogs'
      console.error('Failed to load catalogs:', e)
    } finally {
      isLoading.value = false
    }
  }

  async function createCatalog(name: string, description: string) {
    try {
      const newCatalog = await catalogApi.create({ name, description })
      catalogs.value.push(newCatalog)
      return newCatalog
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create catalog'
      throw e
    }
  }

  async function deleteCatalog(id: string) {
    try {
      await catalogApi.delete(id)
      catalogs.value = catalogs.value.filter(c => c.id !== id)
      if (selectedCatalogId.value === id) {
        selectedCatalogId.value = null
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete catalog'
      throw e
    }
  }

  function selectCatalog(id: string) {
    selectedCatalogId.value = id
  }

  function openCreateDialog() {
    isCreateDialogOpen.value = true
  }

  function closeCreateDialog() {
    isCreateDialogOpen.value = false
  }

  return {
    // State
    catalogs,
    selectedCatalogId,
    isLoading,
    error,
    isCreateDialogOpen,
    
    // Getters
    selectedCatalog,
    catalogCount,
    
    // Actions
    loadCatalogs,
    createCatalog,
    deleteCatalog,
    selectCatalog,
    openCreateDialog,
    closeCreateDialog,
  }
})
```

### Store with Persistence (LocalStorage)

```typescript
// features/tab/stores/tabStore.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Tab } from '../types/tab'

export const useTabStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>([])
  const activeTabId = ref<string | null>(null)

  function createTab(tab: Tab) {
    tabs.value.push(tab)
    activeTabId.value = tab.id
  }

  function closeTab(id: string) {
    tabs.value = tabs.value.filter(t => t.id !== id)
    if (activeTabId.value === id) {
      activeTabId.value = tabs.value[0]?.id ?? null
    }
  }

  return { tabs, activeTabId, createTab, closeTab }
}, {
  persist: {
    key: 'app-tabs',
    storage: localStorage,
  }
})
```

---

## API Service Layer

### Service Structure

**Feature Services** (`features/{feature}/services/`)
- HTTP client for backend API
- Handles request/response transformation
- Error handling
- Type-safe with TypeScript

### Service Template

```typescript
// features/catalog/services/catalogApi.ts
import axios from 'axios'
import type { ApiCatalog, CreateCatalogDto, UpdateCatalogDto } from '../types/catalog'

const API_BASE = '/api'

class CatalogApiService {
  async fetchAll(): Promise<ApiCatalog[]> {
    const response = await axios.get<ApiCatalog[]>(`${API_BASE}/catalogs`)
    return response.data
  }

  async fetchById(id: string): Promise<ApiCatalog> {
    const response = await axios.get<ApiCatalog>(`${API_BASE}/catalogs/${id}`)
    return response.data
  }

  async create(dto: CreateCatalogDto): Promise<ApiCatalog> {
    const response = await axios.post<ApiCatalog>(`${API_BASE}/catalogs`, dto)
    return response.data
  }

  async update(id: string, dto: UpdateCatalogDto): Promise<ApiCatalog> {
    const response = await axios.put<ApiCatalog>(`${API_BASE}/catalogs/${id}`, dto)
    return response.data
  }

  async delete(id: string): Promise<void> {
    await axios.delete(`${API_BASE}/catalogs/${id}`)
  }

  async search(query: string): Promise<ApiCatalog[]> {
    const response = await axios.get<ApiCatalog[]>(
      `${API_BASE}/catalogs/search`,
      { params: { q: query } }
    )
    return response.data
  }
}

// Singleton instance
export const catalogApi = new CatalogApiService()
```

---

## Type Definitions

### Type Organization

**Feature Types** (`features/{feature}/types/`)
- Types specific to feature domain
- DTOs for API communication
- Component prop types

**Shared Types** (`shared/types/`)
- Common types used across features
- Utility types

### Type Template

```typescript
// features/catalog/types/catalog.ts

export interface ApiCatalog {
  id: string
  name: string
  description?: string
  createdAt: Date
  updatedAt: Date
  endpoints: ApiEndpoint[]
  metadata: CatalogMetadata
}

export interface ApiEndpoint {
  id: string
  catalogId: string
  name: string
  method: HttpMethod
  url: string
  description?: string
  headers: Record<string, string>
  queryParams: Parameter[]
  bodyParams: Parameter[]
  pathParams: Parameter[]
  tags: string[]
}

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'HEAD' | 'OPTIONS'

export interface Parameter {
  name: string
  type: 'string' | 'number' | 'boolean' | 'object' | 'array'
  required: boolean
  description?: string
  defaultValue?: any
  enum?: any[]
}

export interface CatalogMetadata {
  source: 'manual' | 'swagger' | 'sandbox'
  swaggerVersion?: string
  originalUrl?: string
}

// DTOs for API calls
export interface CreateCatalogDto {
  name: string
  description?: string
}

export interface UpdateCatalogDto {
  name?: string
  description?: string
}
```

---

## Composables

### Composable Structure

**Feature Composables** (`features/{feature}/composables/`)
- Reusable logic specific to feature
- Encapsulates complex logic
- Returns reactive state and methods

**Shared Composables** (`shared/composables/`)
- Cross-feature utilities
- Examples: debounce, API helpers, toast notifications

### Composable Template

```typescript
// features/catalog/composables/useCatalogSearch.ts
import { ref, computed, type Ref } from 'vue'
import type { ApiCatalog } from '../types/catalog'

export function useCatalogSearch(catalogs: Ref<ApiCatalog[]>) {
  const searchQuery = ref('')

  const filteredCatalogs = computed(() => {
    if (!searchQuery.value) {
      return catalogs.value
    }

    const query = searchQuery.value.toLowerCase()
    
    return catalogs.value.filter(catalog => {
      // Search in catalog name
      if (catalog.name.toLowerCase().includes(query)) {
        return true
      }
      
      // Search in catalog description
      if (catalog.description?.toLowerCase().includes(query)) {
        return true
      }
      
      // Search in endpoint names/URLs
      return catalog.endpoints.some(endpoint => 
        endpoint.name.toLowerCase().includes(query) ||
        endpoint.url.toLowerCase().includes(query)
      )
    })
  })

  return {
    searchQuery,
    filteredCatalogs,
  }
}
```

---

## Router Configuration

### Route Organization

**Feature Routes** (`features/{feature}/router/`)
- Routes specific to feature
- Exported and imported by main router

**Main Router** (`router/index.ts`)
- Combines all feature routes
- Global route configuration

### Feature Routes Template

```typescript
// features/catalog/router/routes.ts
import type { RouteRecordRaw } from 'vue-router'

export const catalogRoutes: RouteRecordRaw[] = [
  {
    path: '/catalogs',
    name: 'catalogs',
    component: () => import('../components/CatalogList.vue'),
    meta: { title: 'Catalogs' }
  },
  {
    path: '/catalogs/:catalogId',
    name: 'catalog-detail',
    component: () => import('../components/CatalogDetail.vue'),
    props: true,
    meta: { title: 'Catalog Detail' }
  },
  {
    path: '/catalogs/:catalogId/endpoints/:endpointId',
    name: 'endpoint-detail',
    component: () => import('../components/EndpointDetail.vue'),
    props: true,
    meta: { title: 'Endpoint Detail' }
  },
]
```

### Main Router Template

```typescript
// router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { catalogRoutes } from '@/features/catalog/router/routes'
import { requestRoutes } from '@/features/request/router/routes'
import { invocationRoutes } from '@/features/invocation/router/routes'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/App.vue'),
    },
    ...catalogRoutes,
    ...requestRoutes,
    ...invocationRoutes,
  ],
})

export default router
```

---

## Cross-Feature Communication

### Option 1: Composables (Recommended for Simple Cases)

```typescript
// features/request/composables/useRequestExecution.ts
import { useCatalogStore } from '@/features/catalog/stores/catalogStore'
import { useInvocationStore } from '@/features/invocation/stores/invocationStore'

export function useRequestExecution() {
  const catalogStore = useCatalogStore()
  const invocationStore = useInvocationStore()

  async function executeFromCatalog(endpointId: string) {
    const endpoint = catalogStore.findEndpoint(endpointId)
    if (!endpoint) return

    const result = await executeRequest(endpoint)
    
    // Capture in flight recorder
    invocationStore.capture({
      url: endpoint.url,
      method: endpoint.method,
      status: result.status,
      response: result.body,
    })
  }

  return { executeFromCatalog }
}
```

### Option 2: Event Bus (For Complex Cases)

```typescript
// shared/composables/useEventBus.ts
import { ref } from 'vue'

type EventHandler = (...args: any[]) => void

class EventBus {
  private events = new Map<string, EventHandler[]>()

  on(event: string, handler: EventHandler) {
    if (!this.events.has(event)) {
      this.events.set(event, [])
    }
    this.events.get(event)!.push(handler)
  }

  off(event: string, handler: EventHandler) {
    const handlers = this.events.get(event)
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  emit(event: string, ...args: any[]) {
    const handlers = this.events.get(event)
    if (handlers) {
      handlers.forEach(handler => handler(...args))
    }
  }
}

const eventBus = new EventBus()

export function useEventBus() {
  return eventBus
}
```

Usage:
```typescript
// Emit event
const eventBus = useEventBus()
eventBus.emit('request:executed', { id: '123', status: 200 })

// Listen to event
onMounted(() => {
  eventBus.on('request:executed', handleRequestExecuted)
})

onUnmounted(() => {
  eventBus.off('request:executed', handleRequestExecuted)
})
```

---

## Feature Flag Pattern

```typescript
// shared/stores/featureFlagsStore.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useFeatureFlagsStore = defineStore('featureFlags', () => {
  const testAutomation = ref(false)

  function enableTestAutomation() {
    testAutomation.value = true
  }

  function disableTestAutomation() {
    testAutomation.value = false
  }

  return {
    testAutomation,
    enableTestAutomation,
    disableTestAutomation,
  }
}, {
  persist: {
    key: 'app-feature-flags',
    storage: localStorage,
  }
})
```

Usage in components:
```vue
<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useFeatureFlagsStore } from '@/shared/stores/featureFlagsStore'

const featureFlagsStore = useFeatureFlagsStore()
const { testAutomation } = storeToRefs(featureFlagsStore)
</script>

<template>
  <div>
    <!-- Only show if feature is enabled -->
    <div v-if="testAutomation" class="scenario-section">
      <h2>Test Scenarios</h2>
      <!-- Scenario components -->
    </div>
  </div>
</template>
```

---

## Styling Guidelines

### Use Tailwind CSS Utility Classes

```vue
<template>
  <div class="flex flex-col h-full bg-gray-50">
    <div class="flex justify-between items-center p-4 bg-white border-b">
      <h2 class="text-xl font-semibold">Catalogs</h2>
      <button class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
        Add Catalog
      </button>
    </div>
  </div>
</template>
```

### Scoped Styles for Component-Specific CSS

```vue
<style scoped>
/* Only when Tailwind classes aren't sufficient */
.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: #888 #f1f1f1;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
}
</style>
```

---

## LLM-Friendly Patterns

### 1. Explicit Imports
```typescript
// Good - explicit imports
import { useCatalogStore } from '@/features/catalog/stores/catalogStore'
import { BaseButton } from '@/shared/components/BaseButton.vue'

// Bad - barrel imports (harder for LLMs to trace)
import { useCatalogStore, useRequestStore } from '@/stores'
```

### 2. Self-Documenting Names
```typescript
// Good
const handleCreateCatalog = () => { ... }
const isCreateDialogOpen = ref(false)

// Bad
const onCreate = () => { ... }
const open = ref(false)
```

### 3. Single Responsibility Components
```vue
<!-- Good - focused component -->
<!-- CatalogSidebar.vue - handles sidebar only -->

<!-- Bad - god component -->
<!-- CatalogView.vue - handles sidebar, detail, dialog, etc. -->
```

### 4. Consistent File Naming
- Components: PascalCase (e.g., `CatalogSidebar.vue`)
- Stores: camelCase + Store (e.g., `catalogStore.ts`)
- Services: camelCase + Api (e.g., `catalogApi.ts`)
- Types: camelCase (e.g., `catalog.ts`)
- Composables: camelCase with "use" prefix (e.g., `useCatalogSearch.ts`)

---

## Testing Strategy

### Component Tests
```typescript
// features/catalog/components/__tests__/CatalogSidebar.test.ts
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import CatalogSidebar from '../CatalogSidebar.vue'

describe('CatalogSidebar', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders catalog list', () => {
    const wrapper = mount(CatalogSidebar)
    expect(wrapper.find('.catalog-sidebar').exists()).toBe(true)
  })

  it('filters catalogs by search query', async () => {
    const wrapper = mount(CatalogSidebar)
    const searchInput = wrapper.find('input')
    
    await searchInput.setValue('API')
    
    // Assertions
  })
})
```

### Store Tests
```typescript
// features/catalog/stores/__tests__/catalogStore.test.ts
import { setActivePinia, createPinia } from 'pinia'
import { useCatalogStore } from '../catalogStore'

describe('catalogStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('loads catalogs', async () => {
    const store = useCatalogStore()
    
    await store.loadCatalogs()
    
    expect(store.catalogs.length).toBeGreaterThan(0)
  })
})
```

---

## Common Mistakes to Avoid

### ❌ Cross-feature component imports
```typescript
// BAD - request feature importing catalog component
import CatalogSidebar from '@/features/catalog/components/CatalogSidebar.vue'
```

### ❌ Shared component importing feature components
```typescript
// BAD - shared component importing feature-specific code
import { useCatalogStore } from '@/features/catalog/stores/catalogStore'
```

### ❌ Mixing business logic in components
```vue
<!-- BAD - API call in component -->
<script setup>
import axios from 'axios'

const loadData = async () => {
  const response = await axios.get('/api/catalogs')
  catalogs.value = response.data
}
</script>

<!-- GOOD - use store or service -->
<script setup>
import { useCatalogStore } from '../stores/catalogStore'

const catalogStore = useCatalogStore()
const loadData = () => catalogStore.loadCatalogs()
</script>
```

---

## Quick Reference

### File Location Decision Tree

```
Is it UI?
├─ Yes → Component
│   ├─ Feature-specific? → features/{feature}/components/
│   └─ Reusable? → shared/components/
│
├─ Is it state?
│   ├─ Feature-specific? → features/{feature}/stores/
│   └─ Cross-feature? → shared/stores/
│
├─ Is it API call?
│   └─ features/{feature}/services/
│
├─ Is it type definition?
│   ├─ Feature-specific? → features/{feature}/types/
│   └─ Common? → shared/types/
│
├─ Is it reusable logic?
│   ├─ Feature-specific? → features/{feature}/composables/
│   └─ Common? → shared/composables/
│
└─ Is it routing?
    └─ features/{feature}/router/
```

---

## Summary

1. **Feature-based** organization mirrors backend
2. **One feature = one directory** for LLM context efficiency
3. **Pinia stores** for state management, collocated with features
4. **API services** abstract HTTP calls from components
5. **Shared code** is minimal - duplicate when in doubt
6. **Explicit imports** for LLM traceability
7. **Consistent naming** for discoverability
8. **Single responsibility** for components and stores