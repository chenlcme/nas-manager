import { createContext } from 'preact';
import { useContext, useState, JSX } from 'preact/hooks';

// SelectionContextType - 选择状态上下文类型
interface SelectionContextType {
  selected: Set<number>;
  toggle: (id: number) => void;
  selectAll: (ids: number[]) => void;
  clear: () => void;
  isSelected: (id: number) => boolean;
  count: number;
}

const SelectionContext = createContext<SelectionContextType | null>(null);

export function SelectionProvider({ children }: { children: JSX.Element }) {
  const [selected, setSelected] = useState<Set<number>>(new Set());

  const toggle = (id: number) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const selectAll = (ids: number[]) => {
    setSelected(new Set(ids));
  };

  const clear = () => {
    setSelected(new Set());
  };

  const isSelected = (id: number) => selected.has(id);

  const value: SelectionContextType = {
    selected,
    toggle,
    selectAll,
    clear,
    isSelected,
    count: selected.size,
  };

  return (
    <SelectionContext.Provider value={value}>
      {children}
    </SelectionContext.Provider>
  );
}

export function useSelection() {
  const context = useContext(SelectionContext);
  if (!context) {
    throw new Error('useSelection must be used within a SelectionProvider');
  }
  return context;
}
