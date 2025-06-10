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
  const dialogRef = React.useRef<HTMLDialogElement>(null);
  
  const handleOpen = () => {
    dialogRef.current?.showModal();
  };
  
  const handleClose = () => {
    dialogRef.current?.close();
  };
  
  const handleBackdropClick = (e: React.MouseEvent<HTMLDialogElement>) => {
    if (e.target === e.currentTarget) {
      handleClose();
    }
  };
  
  return (
    <div className="flex items-center justify-center p-4">
      <Button 
        variant="outline" 
        onClick={handleOpen}
        className="mb-4"
      >
        {buttonText}
      </Button>
      
      <dialog 
        ref={dialogRef}
        className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center backdrop:bg-black/50 backdrop:backdrop-blur-sm"
        onClick={handleBackdropClick}
        aria-labelledby="dialog-title"
      >
        <div 
          className="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full mx-4"
          onClick={(e) => e.stopPropagation()}
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
      </dialog>
    </div>
  );
}