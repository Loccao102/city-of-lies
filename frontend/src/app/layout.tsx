import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "City of Lies — Multi-Agent Investigation",
  description:
    "Real-time investigation game where you must uncover the ground truth before misinformation consumes the city.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="vi" className="dark">
      <body className="bg-background text-foreground antialiased min-h-screen">
        {children}
      </body>
    </html>
  );
}
