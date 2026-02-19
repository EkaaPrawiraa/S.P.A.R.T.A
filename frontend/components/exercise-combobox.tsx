"use client";

import { useMemo, useState } from "react";
import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";

type ExerciseOption = {
  id: string;
  name: string;
};

export function ExerciseCombobox({
  id,
  value,
  onValueChange,
  options,
  placeholder = "Select an exercise",
  disabled,
  className,
  searchPlaceholder = "Search exercise...",
}: {
  id?: string;
  value: string;
  onValueChange: (nextId: string) => void;
  options: ExerciseOption[];
  placeholder?: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  className?: string;
}) {
  const [open, setOpen] = useState(false);

  const sorted = useMemo(() => {
    return [...(options || [])].sort((a, b) =>
      (a.name || "").localeCompare(b.name || "", undefined, {
        sensitivity: "base",
      }),
    );
  }, [options]);

  const selected = useMemo(() => {
    if (!value) return undefined;
    return (options || []).find((o) => o.id === value);
  }, [options, value]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          id={id}
          type="button"
          variant="outline"
          role="combobox"
          aria-expanded={open}
          disabled={disabled}
          className={cn(
            "w-full justify-between",
            !selected && "text-muted-foreground",
            className,
          )}
        >
          <span className="min-w-0 truncate">
            {selected?.name || placeholder}
          </span>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>

      <PopoverContent
        align="start"
        className="w-[--radix-popover-trigger-width] p-0"
      >
        <Command>
          <CommandInput placeholder={searchPlaceholder} />
          <CommandList>
            <CommandEmpty>No exercise found.</CommandEmpty>
            <CommandGroup>
              {sorted.map((opt) => {
                const isSelected = opt.id === value;
                return (
                  <CommandItem
                    key={opt.id}
                    value={opt.name}
                    onSelect={() => {
                      onValueChange(opt.id);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn(
                        "h-4 w-4",
                        isSelected ? "opacity-100" : "opacity-0",
                      )}
                    />
                    <span className="min-w-0 truncate">{opt.name}</span>
                  </CommandItem>
                );
              })}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
