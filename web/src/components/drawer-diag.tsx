"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";

interface DrawerDialogProps {
  title?: string;
  content?: React.ReactNode;
  buttonText?: string;
}

export function DrawerDialog({ 
  title = "Demo Component",
  content = "This is a simple demo component! 🎉",
  buttonText = "Toggle Demo Component"
}: DrawerDialogProps) {
  const [isOpen, setIsOpen] = React.useState(false);
  
  const handleClose = () => setIsOpen(false);
  
  // ESC key handler
  React.useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        handleClose();
      }
    };
    
    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      return () => document.removeEventListener('keydown', handleEscape);
    }
  }, [isOpen]);
  
  return (
    <div className="flex items-center justify-center p-4">
      <Button 
        variant="outline" 
        onClick={() => setIsOpen(!isOpen)}
        className="mb-4"
      >
        {buttonText}
      </Button>
      
      {isOpen && (
        <div 
          className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center"
          onClick={handleClose}
          onKeyDown={(e) => e.key === 'Enter' && handleClose()}
          role="dialog"
          aria-modal="true"
          aria-labelledby="dialog-title"
          tabIndex={0}
        >
          <div 
            className="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full mx-4"
            onClick={(e) => e.stopPropagation()}
            onKeyDown={(e) => e.stopPropagation()}
            tabIndex={0}
          >
            <h2 id="dialog-title" className="text-lg font-semibold mb-2">
              {title}
            </h2>
            <div className="text-gray-600 mb-4">
              {content}
            </div>
            <Button 
              onClick={handleClose}
              className="w-full"
            >
              Close
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}