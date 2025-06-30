"use client"

import { useState } from "react"
import { MapPinIcon, MenuIcon, SearchIcon, HelpCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between">
        <div className="flex items-center">
          <DropdownMenu open={mobileMenuOpen} onOpenChange={setMobileMenuOpen}>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden mr-4" aria-label="モバイルメニューを開く" aria-expanded={mobileMenuOpen}>
                <MenuIcon className="h-5 w-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-56">
              <DropdownMenuItem aria-label="スポットを探す">
                <SearchIcon className="h-4 w-4 mr-2" />
                スポットを探す
              </DropdownMenuItem>
              <DropdownMenuItem aria-label="レビューを書く">
                レビューを書く
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          
          <div className="hidden md:flex items-center space-x-2">
            <Button variant="ghost" size="sm">
              <SearchIcon className="h-4 w-4 mr-2" />
              スポットを探す
            </Button>
            <Button variant="ghost" size="sm">
              レビューを書く
            </Button>
          </div>
        </div>
        
        <div className="absolute left-1/2 transform -translate-x-1/2 flex items-center space-x-2">
          <MapPinIcon className="h-6 w-6 text-primary" aria-hidden="true" />
          <h1 className="text-xl font-bold">Bocchi The Map</h1>
        </div>
        
        <div className="flex items-center space-x-4">
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="sm" aria-label="ヘルプを表示">
                <HelpCircle className="h-4 w-4" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80">
              <div className="space-y-4">
                <h3 className="font-medium text-sm text-center">アプリの特徴</h3>
                <div className="space-y-3">
                  <div>
                    <h4 className="font-semibold text-sm mb-1">🎯 簡単検索</h4>
                    <p className="text-xs text-muted-foreground">
                      カテゴリーや場所から、あなたにぴったりのスポットを見つけましょう
                    </p>
                  </div>
                  <div>
                    <h4 className="font-semibold text-sm mb-1">💬 リアルな口コミ</h4>
                    <p className="text-xs text-muted-foreground">
                      実際に訪れた人のレビューで、一人でも入りやすいお店がわかります
                    </p>
                  </div>
                  <div>
                    <h4 className="font-semibold text-sm mb-1">📍 マップ表示</h4>
                    <p className="text-xs text-muted-foreground">
                      現在地から近いおひとりさまスポットを地図上で確認できます
                    </p>
                  </div>
                </div>
              </div>
            </PopoverContent>
          </Popover>
          <Button variant="default" size="sm">
            ログイン
          </Button>
        </div>
      </div>
    </header>
  )
}