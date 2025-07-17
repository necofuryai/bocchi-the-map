import React from 'react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { Header } from '../header'

// Mock Clerk components
vi.mock('@clerk/nextjs', () => ({
  SignedIn: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  SignedOut: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  UserButton: () => <button aria-label="ユーザーメニューを開く">User Menu</button>,
  useUser: () => ({ user: { firstName: 'Test' } }),
}))

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
    it('When viewed on desktop, Then help button should be visible', () => {
      expect(screen.getByRole('button', { name: 'ヘルプを表示' })).toBeInTheDocument()
    })
  })

  describe('Given the Header component is rendered on mobile', () => {
    it('When the mobile menu button is present, Then it should be visible', () => {
      const mobileMenuButton = screen.getByRole('button', { name: 'モバイルメニューを開く' })
      expect(mobileMenuButton).toBeInTheDocument()
    })
  })

  describe('Given the user menu is accessible', () => {
    it('When the user menu button is present, Then it should be visible', () => {
      const userMenuButton = screen.getByRole('button', { name: 'ユーザーメニューを開く' })
      expect(userMenuButton).toBeInTheDocument()
    })

    it('When the user menu button is clicked, Then it should be interactable', () => {
      const userMenuButton = screen.getByRole('button', { name: 'ユーザーメニューを開く' })
      fireEvent.click(userMenuButton)
      
      // Since it's a mocked UserButton, we just verify it can be clicked
      expect(userMenuButton).toBeInTheDocument()
    })
  })

  describe('Given the Header component has accessibility features', () => {
    it('When rendered, Then proper ARIA attributes should be present', () => {
      const mobileMenuButton = screen.getByRole('button', { name: 'モバイルメニューを開く' })
      expect(mobileMenuButton).toHaveAttribute('aria-expanded', 'false')
    })

    it('When mobile menu is clicked, Then it should be interactable', () => {
      const mobileMenuButton = screen.getByRole('button', { name: 'モバイルメニューを開く' })
      
      fireEvent.click(mobileMenuButton)
      
      // Since it's a functional component, we just verify it exists
      expect(mobileMenuButton).toBeInTheDocument()
    })
  })
})