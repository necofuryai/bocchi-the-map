import React from 'react'
import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { Header } from '../header'

describe('Header Component', () => {
  beforeEach(() => {
    render(<Header />)
  })

  describe('Given the Header component is rendered', () => {
    it('When the component loads, Then the application title should be displayed', () => {
      expect(screen.getByText('Bocchi The Map')).toBeInTheDocument()
    })

    it('When the component loads, Then the map pin icon should be visible', () => {
      const mapPinIcon = screen.getByRole('heading', { name: 'Bocchi The Map' }).parentNode?.querySelector('svg')
      expect(mapPinIcon).toBeInTheDocument()
    })

    it('When the component loads, Then the user menu button should be visible', () => {
      const userMenuButton = screen.getByRole('button', { name: 'ユーザーメニューを開く' })
      expect(userMenuButton).toBeInTheDocument()
    })
  })

  describe('Given the Header component is rendered on desktop', () => {
    it('When viewed on desktop, Then navigation buttons should be visible', () => {
      expect(screen.getByText('スポットを探す')).toBeInTheDocument()
      expect(screen.getByText('レビューを書く')).toBeInTheDocument()
    })
  })

  describe('Given the Header component is rendered on mobile', () => {
    it('When the mobile menu button is clicked, Then the mobile menu should open', () => {
      const mobileMenuButton = screen.getByRole('button', { name: 'モバイルメニューを開く' })
      
      fireEvent.click(mobileMenuButton)
      
      // Check that aria-expanded is updated for accessibility
      expect(mobileMenuButton).toHaveAttribute('aria-expanded', 'true')
      
      // Check if mobile menu items are visible
      const mobileSearchButton = screen.getAllByText('スポットを探す').find(el => 
        el.closest('[role="menuitem"]')
      )
      expect(mobileSearchButton).toBeInTheDocument()
    })
  })

  describe('Given the user menu is accessible', () => {
    it('When the user menu button is clicked, Then the user menu should open', () => {
      const userMenuButton = screen.getByRole('button', { name: 'ユーザーメニューを開く' })
      
      fireEvent.click(userMenuButton)
      
      expect(screen.getByText('マイアカウント')).toBeInTheDocument()
      expect(screen.getByText('プロフィール')).toBeInTheDocument()
      expect(screen.getByText('レビュー履歴')).toBeInTheDocument()
      expect(screen.getByText('お気に入り')).toBeInTheDocument()
      expect(screen.getByText('ログアウト')).toBeInTheDocument()
    })

    it('When the user menu is open, Then all menu items should have proper aria labels', () => {
      const userMenuButton = screen.getByRole('button', { name: 'ユーザーメニューを開く' })
      fireEvent.click(userMenuButton)
      
      expect(screen.getByRole('menuitem', { name: 'プロフィールページを表示' })).toBeInTheDocument()
      expect(screen.getByRole('menuitem', { name: 'レビュー履歴ページを表示' })).toBeInTheDocument()
      expect(screen.getByRole('menuitem', { name: 'お気に入りスポット一覧ページを表示' })).toBeInTheDocument()
      expect(screen.getByRole('menuitem', { name: 'アカウントからログアウトする' })).toBeInTheDocument()
    })
  })

  describe('Given the Header component has accessibility features', () => {
    it('When rendered, Then proper ARIA attributes should be present', () => {
      const userMenuButton = screen.getByRole('button', { name: 'ユーザーメニューを開く' })
      expect(userMenuButton).toHaveAttribute('aria-expanded', 'false')
      
      const mobileMenuButton = screen.getByRole('button', { name: 'モバイルメニューを開く' })
      expect(mobileMenuButton).toHaveAttribute('aria-expanded', 'false')
    })

    it('When menus are opened, Then aria-expanded should be updated', () => {
      const userMenuButton = screen.getByRole('button', { name: 'ユーザーメニューを開く' })
      
      fireEvent.click(userMenuButton)
      
      expect(userMenuButton).toHaveAttribute('aria-expanded', 'true')
    })
  })
})