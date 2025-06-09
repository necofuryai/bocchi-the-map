"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { CoffeeIcon, UtensilsIcon, ShoppingBagIcon, BookOpenIcon, MoreHorizontalIcon } from "lucide-react"

const categories = [
  { icon: CoffeeIcon, name: "カフェ", count: 156 },
  { icon: UtensilsIcon, name: "レストラン", count: 234 },
  { icon: ShoppingBagIcon, name: "ショッピング", count: 89 },
  { icon: BookOpenIcon, name: "書店・図書館", count: 67 },
  { icon: MoreHorizontalIcon, name: "その他", count: 123 },
]

export function Sidebar() {
  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">カテゴリー</CardTitle>
          <CardDescription>
            おひとりさまに人気のスポットを探そう
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-2">
          {categories.map((category) => (
            <Button
              key={category.name}
              variant="ghost"
              className="w-full justify-start"
            >
              <category.icon className="mr-2 h-4 w-4" />
              {category.name}
              <span className="ml-auto text-muted-foreground">
                {category.count}
              </span>
            </Button>
          ))}
        </CardContent>
      </Card>
      
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">最近のレビュー</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="text-sm">
            <p className="font-medium">渋谷の隠れ家カフェ</p>
            <p className="text-muted-foreground">2時間前</p>
          </div>
          <div className="text-sm">
            <p className="font-medium">新宿の静かな図書館</p>
            <p className="text-muted-foreground">5時間前</p>
          </div>
          <div className="text-sm">
            <p className="font-medium">原宿の一人焼肉店</p>
            <p className="text-muted-foreground">1日前</p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}