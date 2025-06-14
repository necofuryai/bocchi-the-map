import React from 'react'
import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { Button } from '../button'

describe('Button Component', () => {
  describe('Given the Button component', () => {
    it('When rendered with default props, Then it should display the button text', () => {
      render(<Button>Click me</Button>)
      
      expect(screen.getByRole('button', { name: 'Click me' })).toBeInTheDocument()
    })

    it('When rendered with default type, Then it should have type="button"', () => {
      render(<Button>Default Button</Button>)
      
      expect(screen.getByRole('button')).toHaveAttribute('type', 'button')
    })

    it('When type prop is provided, Then it should use the specified type', () => {
      render(<Button type="submit">Submit Button</Button>)
      
      expect(screen.getByRole('button')).toHaveAttribute('type', 'submit')
    })
  })

  describe('Given the Button component with different variants', () => {
    it('When variant is "default", Then it should have primary background styles', () => {
      render(<Button variant="default">Default Button</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('bg-primary', 'text-primary-foreground', 'shadow')
    })

    it('When variant is "destructive", Then it should have destructive background styles', () => {
      render(<Button variant="destructive">Delete Button</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('bg-destructive', 'text-destructive-foreground', 'shadow-sm')
    })

    it('When variant is "outline", Then it should have outline border styles', () => {
      render(<Button variant="outline">Outline Button</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('border', 'border-input', 'bg-background', 'shadow-sm')
    })

    it('When variant is "ghost", Then it should have ghost hover styles', () => {
      render(<Button variant="ghost">Ghost Button</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('hover:bg-accent', 'hover:text-accent-foreground')
    })
  })

  describe('Given the Button component with different sizes', () => {
    it('When size is "default", Then it should have default size styles', () => {
      render(<Button size="default">Default Size</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('h-9', 'px-4', 'py-2')
    })

    it('When size is "sm", Then it should have small size styles', () => {
      render(<Button size="sm">Small Button</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('h-8', 'rounded-md', 'px-3', 'text-xs')
    })

    it('When size is "lg", Then it should have large size styles', () => {
      render(<Button size="lg">Large Button</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('h-10', 'rounded-md', 'px-8')
    })

    it('When size is "icon", Then it should have icon size styles', () => {
      render(<Button size="icon">Icon</Button>)
      
      const button = screen.getByRole('button')
      expect(button).toHaveClass('h-9', 'w-9')
    })
  })

  describe('Given the Button component with event handlers', () => {
    it('When clicked, Then it should call the onClick handler', () => {
      const handleClick = vi.fn()
      render(<Button onClick={handleClick}>Clickable Button</Button>)
      
      fireEvent.click(screen.getByRole('button'))
      
      expect(handleClick).toHaveBeenCalledTimes(1)
    })

    it('When disabled, Then onClick should not be called', () => {
      const handleClick = vi.fn()
      render(<Button onClick={handleClick} disabled>Disabled Button</Button>)
      
      fireEvent.click(screen.getByRole('button'))
      
      expect(handleClick).not.toHaveBeenCalled()
    })
  })

  describe('Given the Button component with accessibility features', () => {
    it('When disabled, Then it should have disabled attribute', () => {
      render(<Button disabled>Disabled Button</Button>)
      
      expect(screen.getByRole('button')).toBeDisabled()
    })

    it('When aria-label is provided, Then it should have the aria-label', () => {
      render(<Button aria-label="Close dialog">×</Button>)
      
      expect(screen.getByRole('button', { name: 'Close dialog' })).toBeInTheDocument()
    })

    it('When custom className is provided, Then it should merge with default classes', () => {
      render(<Button className="custom-class">Custom Button</Button>)
      
      const button = screen.getByRole('button')
      const classNames = button.className.split(' ')
      expect(classNames).toEqual(expect.arrayContaining(['custom-class', 'inline-flex', 'items-center', 'justify-center']))
    })
  })
})