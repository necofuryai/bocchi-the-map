"use client"

import { MapPinIcon, MenuIcon, SearchIcon, UserIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between">
        <div className="flex items-center">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden mr-4" aria-label="モバイルメニューを開く">
                <MenuIcon className="h-5 w-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-56">
              <DropdownMenuItem>
                <SearchIcon className="h-4 w-4 mr-2" />
                スポットを探す
              </DropdownMenuItem>
              <DropdownMenuItem>
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
          <MapPinIcon className="h-6 w-6 text-primary" />
          <h1 className="text-xl font-bold">Bocchi The Map</h1>
        </div>
        
        <div className="flex items-center space-x-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" aria-label="ユーザーメニューを開く">
                <UserIcon className="h-5 w-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>マイアカウント</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                プロフィール
              </DropdownMenuItem>
              <DropdownMenuItem>
                レビュー履歴
              </DropdownMenuItem>
              <DropdownMenuItem>
                お気に入り
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                ログアウト
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}