'use client'
import "./globals.css";
import Navigation from '@/components/Navigation';


export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html>
      <body>
        <Navigation />
        {children}
      </body>
    </html>
  );
}