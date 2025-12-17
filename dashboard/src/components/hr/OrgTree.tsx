import * as React from 'react';
import { ChevronRight, ChevronDown, Building2, FolderTree } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { OrgUnitSnapshot } from '@/types/hr';

interface OrgTreeProps {
  orgUnits: OrgUnitSnapshot[];
  onSelect?: (orgUnit: OrgUnitSnapshot) => void;
  selectedId?: string;
}

interface TreeNode extends OrgUnitSnapshot {
  children: TreeNode[];
}

function buildTree(orgUnits: OrgUnitSnapshot[]): TreeNode[] {
  const nodeMap = new Map<string, TreeNode>();
  const roots: TreeNode[] = [];

  // Create node map
  orgUnits.forEach((unit) => {
    nodeMap.set(unit.orgUnitId, { ...unit, children: [] });
  });

  // Build tree structure
  orgUnits.forEach((unit) => {
    const node = nodeMap.get(unit.orgUnitId)!;
    if (unit.parentId && nodeMap.has(unit.parentId)) {
      nodeMap.get(unit.parentId)!.children.push(node);
    } else {
      roots.push(node);
    }
  });

  // Sort children alphabetically
  const sortChildren = (nodes: TreeNode[]) => {
    nodes.sort((a, b) => a.name.localeCompare(b.name));
    nodes.forEach((node) => sortChildren(node.children));
  };
  sortChildren(roots);

  return roots;
}

interface TreeNodeItemProps {
  node: TreeNode;
  level: number;
  onSelect?: (orgUnit: OrgUnitSnapshot) => void;
  selectedId?: string;
  expandedIds: Set<string>;
  toggleExpand: (id: string) => void;
}

function TreeNodeItem({ node, level, onSelect, selectedId, expandedIds, toggleExpand }: TreeNodeItemProps) {
  const hasChildren = node.children.length > 0;
  const isExpanded = expandedIds.has(node.orgUnitId);
  const isSelected = selectedId === node.orgUnitId;

  return (
    <div>
      <div
        className={cn(
          'flex items-center gap-2 py-1.5 px-2 rounded-md cursor-pointer transition-colors',
          isSelected ? 'bg-green-100 text-green-800' : 'hover:bg-gray-100'
        )}
        style={{ paddingLeft: `${level * 20 + 8}px` }}
        onClick={() => onSelect?.(node)}
      >
        {hasChildren ? (
          <button
            onClick={(e) => {
              e.stopPropagation();
              toggleExpand(node.orgUnitId);
            }}
            className="p-0.5 hover:bg-gray-200 rounded"
          >
            {isExpanded ? (
              <ChevronDown className="h-4 w-4 text-gray-500" />
            ) : (
              <ChevronRight className="h-4 w-4 text-gray-500" />
            )}
          </button>
        ) : (
          <span className="w-5" />
        )}
        <Building2 className="h-4 w-4 text-green-600 flex-shrink-0" />
        <span className="font-medium truncate">{node.name}</span>
        <span className="text-xs text-gray-500 ml-1">({node.code})</span>
        {node.kind && (
          <span className="text-xs bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded ml-auto">
            {node.kind}
          </span>
        )}
      </div>
      {hasChildren && isExpanded && (
        <div>
          {node.children.map((child) => (
            <TreeNodeItem
              key={child.orgUnitId}
              node={child}
              level={level + 1}
              onSelect={onSelect}
              selectedId={selectedId}
              expandedIds={expandedIds}
              toggleExpand={toggleExpand}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export function OrgTree({ orgUnits, onSelect, selectedId }: OrgTreeProps) {
  const [expandedIds, setExpandedIds] = React.useState<Set<string>>(new Set());

  const tree = React.useMemo(() => buildTree(orgUnits), [orgUnits]);

  // Expand all by default when tree changes
  React.useEffect(() => {
    const allIds = new Set<string>();
    const collectIds = (nodes: TreeNode[]) => {
      nodes.forEach((node) => {
        if (node.children.length > 0) {
          allIds.add(node.orgUnitId);
          collectIds(node.children);
        }
      });
    };
    collectIds(tree);
    setExpandedIds(allIds);
  }, [tree]);

  const toggleExpand = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const expandAll = () => {
    const allIds = new Set<string>();
    const collectIds = (nodes: TreeNode[]) => {
      nodes.forEach((node) => {
        if (node.children.length > 0) {
          allIds.add(node.orgUnitId);
          collectIds(node.children);
        }
      });
    };
    collectIds(tree);
    setExpandedIds(allIds);
  };

  const collapseAll = () => {
    setExpandedIds(new Set());
  };

  if (orgUnits.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-gray-500">
        <FolderTree className="h-12 w-12 mb-4 text-gray-300" />
        <p>No organizational units found</p>
        <p className="text-sm">Create org units to see the hierarchy</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-sm text-gray-600">
          <FolderTree className="h-4 w-4" />
          <span>{orgUnits.length} org units</span>
        </div>
        <div className="flex gap-2">
          <button
            onClick={expandAll}
            className="text-xs text-gray-500 hover:text-gray-700 px-2 py-1 hover:bg-gray-100 rounded"
          >
            Expand All
          </button>
          <button
            onClick={collapseAll}
            className="text-xs text-gray-500 hover:text-gray-700 px-2 py-1 hover:bg-gray-100 rounded"
          >
            Collapse All
          </button>
        </div>
      </div>
      <div className="border rounded-lg p-2 bg-white">
        {tree.map((node) => (
          <TreeNodeItem
            key={node.orgUnitId}
            node={node}
            level={0}
            onSelect={onSelect}
            selectedId={selectedId}
            expandedIds={expandedIds}
            toggleExpand={toggleExpand}
          />
        ))}
      </div>
    </div>
  );
}
