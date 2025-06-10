"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";

export function DrawerDialogDemo() {
  const [isOpen, setIsOpen] = React.useState(false);
  
  return (
    <div className="flex items-center justify-center p-4">
      <Button 
        variant="outline" 
        onClick={() => setIsOpen(!isOpen)}
        className="mb-4"
      >
        Toggle Demo Component
      </Button>
      
      {isOpen && (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center">
          <div className="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full mx-4">
            <h2 className="text-lg font-semibold mb-2">Demo Component</h2>
            <p className="text-gray-600 mb-4">
              This is a simple demo component! 🎉
            </p>
            <Button 
              onClick={() => setIsOpen(false)}
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