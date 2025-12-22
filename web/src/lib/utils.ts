import type { ClassValue } from "clsx"
import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Tree node type with children
 */
export type TreeNode<T> = T & {
  children?: TreeNode<T>[]
}

/**
 * Convert flat array to tree structure
 * @param items Flat array of items with id and parentId
 * @param options Configuration options
 * @returns Array of root nodes with children
 */
export function arrayToTree<T extends { id: number, parentId: number | null }>(
  items: T[],
  options?: {
    rootParentId?: number | null
    sortBy?: (a: T, b: T) => number
  }
): TreeNode<T>[] {
  const { rootParentId = null, sortBy } = options || {}

  // Create a map for quick lookup
  const itemMap = new Map<number, TreeNode<T>>()
  const result: TreeNode<T>[] = []

  // First pass: create map of all items
  items.forEach((item) => {
    itemMap.set(item.id, { ...item, children: [] })
  })

  // Second pass: build tree structure
  items.forEach((item) => {
    const node = itemMap.get(item.id)!

    if (item.parentId === rootParentId) {
      // Root node
      result.push(node)
    } else if (item.parentId !== null) {
      // Child node
      const parent = itemMap.get(item.parentId)
      if (parent) {
        if (!parent.children) {
          parent.children = []
        }
        parent.children.push(node)
      } else {
        // Parent not found, treat as root
        result.push(node)
      }
    }
  })

  // Sort if sortBy function provided
  if (sortBy) {
    const sortRecursive = (nodes: TreeNode<T>[]) => {
      nodes.sort(sortBy)
      nodes.forEach((node) => {
        if (node.children && node.children.length > 0) {
          sortRecursive(node.children)
        }
      })
    }
    sortRecursive(result)
  }

  return result
}
