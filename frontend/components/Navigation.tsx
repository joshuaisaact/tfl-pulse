import Link from 'next/link';
import { usePathname } from 'next/navigation';

export default function Navigation() {
    const pathname = usePathname();

    return (
        <nav className="bg-[#009FE0] shadow-lg">
            <div className="max-w-7xl mx-auto px-4">
                <div className="flex space-x-4 h-16 items-center">
                    <Link
                        href="/"
                        className={`px-3 py-2 rounded-md text-sm font-medium text-white ${
                            pathname === '/' ? 'bg-[#0077AA]' : 'hover:bg-[#0077AA]'
                        }`}
                    >
                        Train List
                    </Link>
                    <Link
                        href="/map"
                        className={`px-3 py-2 rounded-md text-sm font-medium text-white ${
                            pathname === '/map' ? 'bg-[#0077AA]' : 'hover:bg-[#0077AA]'
                        }`}
                    >
                        Line Map
                    </Link>
                </div>
            </div>
        </nav>
    );
}