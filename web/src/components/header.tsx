"use client"

import React, { useState } from "react"
import { MapPinIcon, Bars3Icon, QuestionMarkCircleIcon } from "@heroicons/react/24/outline"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { useAuth } from "@/components/providers/user-provider"
import { UserProfile } from "@/components/auth/user-profile"

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const { user, loading } = useAuth()

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between">
        <div className="flex items-center">
          <DropdownMenu open={mobileMenuOpen} onOpenChange={setMobileMenuOpen}>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden mr-4" aria-label="モバイルメニューを開く" aria-expanded={mobileMenuOpen}>
                <Bars3Icon className="h-5 w-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-56">
              {/* MVPでは検索機能・レビュー機能は未実装 */}
            </DropdownMenuContent>
          </DropdownMenu>
          
          <div className="hidden md:flex items-center space-x-2">
            {/* MVPでは検索機能・レビュー機能は未実装 */}
          </div>
        </div>
        
        <div className="absolute left-1/2 transform -translate-x-1/2 flex items-center space-x-2">
          <Link href="/" className="flex items-center space-x-2">
            <MapPinIcon className="h-6 w-6 text-primary" aria-hidden="true" />
            <h1 className="text-xl font-bold">Bocchi The Map</h1>
          </Link>
        </div>
        
        <div className="flex items-center space-x-4">
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="sm" aria-label="ヘルプを表示">
                <QuestionMarkCircleIcon className="h-4 w-4" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80">
              <div className="space-y-4">
                <h3 className="font-medium text-sm text-center">アプリの特徴</h3>
                <div className="space-y-3">
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
          
          {/* Authentication UI */}
          {loading ? (
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
          ) : user ? (
            <UserProfile />
          ) : (
            <div className="flex items-center">
              <Link href="/auth/login">
                <Button variant="default" size="sm">ログイン / 新規登録</Button>
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}