"use client"

import React, { useState } from "react"
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
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [userMenuOpen, setUserMenuOpen] = useState(false)
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
          <MapPinIcon className="h-6 w-6 text-primary" />
          <h1 className="text-xl font-bold">Bocchi The Map</h1>
        </div>
        
        <div className="flex items-center space-x-4">
          <DropdownMenu open={userMenuOpen} onOpenChange={setUserMenuOpen}>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" aria-label="ユーザーメニューを開く" aria-expanded={userMenuOpen}>
                <UserIcon className="h-5 w-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>マイアカウント</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem aria-label="プロフィールページを表示">
                プロフィール
              </DropdownMenuItem>
              <DropdownMenuItem aria-label="レビュー履歴ページを表示">
                レビュー履歴
              </DropdownMenuItem>
              <DropdownMenuItem aria-label="お気に入りスポット一覧ページを表示">
                お気に入り
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem aria-label="アカウントからログアウトする">
                ログアウト
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}