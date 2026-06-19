import { useState, useEffect, createContext, useContext } from 'react'
import { BrowserRouter, Routes, Route, Link, Outlet, useNavigate, useLocation } from 'react-router-dom'
import { QueryClient, QueryClientProvider, useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import { Moon, Sun, Loader2, Settings, Plus, LogOut } from 'lucide-react'

// 导入具体页面视图组件
import HomeView from './views/HomeView'
import BlogView from './views/BlogView'
import LinksView from './views/LinksView'
import LoginView from './views/LoginView'
import PostDetailView from './views/PostDetailView'
import NewPostView from './views/NewPostView'
import SettingsView from './views/SettingsView'
import EditPostView from './views/EditPostView'
import EditHomeView from './views/EditHomeView'
import EditLinksView from './views/EditLinksView'

// 1. 初始化 React Query 客户端
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60 * 1000,
      refetchOnWindowFocus: false,
      retry: false,
    },
  },
})

// 2. 统一 Axios 默认配置（支持 Cookie 凭证传递）
axios.defaults.withCredentials = true

// 3. 全局上下文 (Context) 定义
interface AuthContextType {
  user: { username: string } | null
  isLoading: boolean
  siteConfig: Record<string, string>
  refetchAuth: () => void
  refetchConfig: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within an AuthProvider')
  return context
}

// 4. 全局母版布局组件 (Layout)
function Layout() {
  const { user, siteConfig } = useAuth()
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  
  // 主题管理
  const [theme, setTheme] = useState(() => localStorage.getItem('theme') || 'dark')
  
  useEffect(() => {
    const root = window.document.documentElement
    if (theme === 'dark') {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
    localStorage.setItem('theme', theme)
  }, [theme])

  const toggleTheme = () => setTheme(prev => prev === 'dark' ? 'light' : 'dark')

  // 退出登录
  const logoutMutation = useMutation({
    mutationFn: () => axios.post('/api/auth/logout'),
    onSuccess: () => {
      queryClient.setQueryData(['auth-user'], null)
      queryClient.invalidateQueries({ queryKey: ['auth-user'] })
      navigate('/login')
    }
  })

  const blogName = siteConfig.blog_name || "Bitlog"

  return (
    <div className="min-h-screen flex flex-col antialiased bg-background text-foreground transition-colors duration-300">
      <header className="py-6 shrink-0">
        <div className="max-w-3xl mx-auto px-5 flex items-center justify-between flex-wrap gap-5">
          <div className="flex items-center gap-5 flex-wrap">
            <Link to="/" className="text-base font-bold tracking-tight">
              {blogName}
            </Link>
            <nav className="flex gap-5 items-center">
              <Link to="/blog" className="text-blue-600 dark:text-blue-400 hover:underline text-base font-medium">
                Blog
              </Link>
              <Link to="/links" className="text-blue-600 dark:text-blue-400 hover:underline text-base font-medium">
                Links
              </Link>
              {user && (
                <>
                  <span className="w-px h-4 bg-neutral-300 dark:bg-neutral-700" />
                  <Link to="/new" className="text-sm text-neutral-500 hover:text-neutral-900 dark:text-neutral-400 dark:hover:text-neutral-100 transition-colors flex items-center gap-1">
                    <Plus size={14} /> New
                  </Link>
                  <Link to="/settings" className="text-sm text-neutral-500 hover:text-neutral-900 dark:text-neutral-400 dark:hover:text-neutral-100 transition-colors flex items-center gap-1">
                    <Settings size={14} /> Settings
                  </Link>
                  <button 
                    onClick={() => logoutMutation.mutate()}
                    className="text-sm text-neutral-500 hover:text-red-500 dark:text-neutral-400 dark:hover:text-red-400 transition-colors flex items-center gap-1 cursor-pointer bg-transparent border-0"
                    title="退出登录"
                  >
                    <LogOut size={14} />
                  </button>
                </>
              )}
            </nav>
          </div>
          
          <button
            onClick={toggleTheme}
            className="p-2 text-neutral-600 dark:text-neutral-400 hover:bg-neutral-100 dark:hover:bg-neutral-800 rounded-md transition-colors flex items-center justify-center cursor-pointer border-0 bg-transparent"
            aria-label="Toggle Theme"
          >
            {theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
          </button>
        </div>
      </header>

      <main className="flex-grow w-full">
        <Outlet />
      </main>

      <footer className="py-10 text-center text-sm text-neutral-500 mt-auto px-5 shrink-0">
        <p>Powered by <a href="https://github.com/ifndf/limblog" target="_blank" rel="noopener noreferrer" className="hover:underline">Bitlog</a></p>
      </footer>
    </div>
  )
}

// 5. 权限与私密拦截代理 (Guard)
function AuthProvider({ children }: { children: React.ReactNode }) {
  // 获取站点配置
  const { data: siteConfigData, isLoading: isConfigLoading, refetch: refetchConfig } = useQuery({
    queryKey: ['site-config'],
    queryFn: () => axios.get('/api/site-config').then(res => res.data),
  })

  // 获取当前登录用户
  const { data: userData, isLoading: isUserLoading, refetch: refetchAuth } = useQuery({
    queryKey: ['auth-user'],
    queryFn: () => axios.get('/api/auth/me').then(res => res.data).catch(() => null),
  })

  const isLoading = isConfigLoading || isUserLoading
  const siteConfig = siteConfigData || {}

  // 博客可见性隐私拦截（伪装 404）
  const location = useLocation()
  const isPrivateRoute = siteConfig.site_visibility === 'private'
  const isPublicPath = !location.pathname.startsWith('/login') && !location.pathname.startsWith('/settings') && !location.pathname.startsWith('/new') && !location.pathname.startsWith('/edit') && !location.pathname.startsWith('/home') && !location.pathname.startsWith('/links/edit')
  
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#0f0f0f]">
        <Loader2 className="animate-spin text-neutral-500" size={32} />
      </div>
    )
  }

  // 私密隐身状态拦截
  if (isPrivateRoute && isPublicPath && !userData) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center p-5 text-center">
        <h1 className="text-6xl font-black text-stone-700 dark:text-neutral-800">404</h1>
        <p className="text-neutral-500 mt-2">Page Not Found</p>
        <Link to="/login" className="text-blue-500 hover:underline mt-4 text-sm">Admin Entry</Link>
      </div>
    )
  }

  return (
    <AuthContext.Provider value={{ user: userData, isLoading, siteConfig, refetchAuth, refetchConfig }}>
      {children}
    </AuthContext.Provider>
  )
}

// 路由拦截器组件
function RequireAuth({ children }: { children: React.ReactNode }) {
  const { user } = useAuth()
  const navigate = useNavigate()
  
  useEffect(() => {
    if (!user) {
      navigate('/login', { replace: true })
    }
  }, [user, navigate])

  if (!user) return null
  return <>{children}</>
}

// 6. 前端路由映射入口
export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            <Route path="/" element={<Layout />}>
              {/* 公开页面 */}
              <Route index element={<HomeView />} />
              <Route path="blog" element={<BlogView />} />
              <Route path="links" element={<LinksView />} />
              <Route path="login" element={<LoginView />} />
              <Route path=":slug" element={<PostDetailView />} />
              
              {/* 后台管理页面 (受权限中间件保护) */}
              <Route path="new" element={<RequireAuth><NewPostView /></RequireAuth>} />
              <Route path="settings" element={<RequireAuth><SettingsView /></RequireAuth>} />
              <Route path="edit/:slug" element={<RequireAuth><EditPostView /></RequireAuth>} />
              <Route path="home/edit" element={<RequireAuth><EditHomeView /></RequireAuth>} />
              <Route path="links/edit" element={<RequireAuth><EditLinksView /></RequireAuth>} />
              
              {/* 兜底 404 页 */}
              <Route path="*" element={<NotFoundView />} />
            </Route>
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

function NotFoundView() {
  return (
    <div className="max-w-3xl mx-auto px-5 py-20 text-center">
      <h1 className="text-4xl font-extrabold">404</h1>
      <p className="text-neutral-500 mt-2">找不到该页面</p>
      <Link to="/" className="text-blue-500 hover:underline mt-4 inline-block">返回首页</Link>
    </div>
  )
}
