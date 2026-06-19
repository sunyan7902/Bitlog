import { Link } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import remarkBreaks from 'remark-breaks'
import { Mail, Pencil } from 'lucide-react'
import { useAuth } from '../App'

export default function HomeView() {
  const { user, siteConfig } = useAuth()

  const defaultHomeContent = `Welcome to **Bitlog** — a minimalist space built for writing that matters.

Here you can capture technical insights, document your thinking, or jot down ideas worth keeping.

Browse the [Blog](/blog) for all posts, or head to [Admin](/login) to manage your content.

Less noise. More words.`

  const homeContent = siteConfig.home_content?.trim() || defaultHomeContent

  return (
    <div className="max-w-3xl mx-auto px-5 space-y-10 animate-fade-in relative">
      {user && (
        <div className="flex justify-end mb-4 absolute top-0 right-5">
          <Link
            to="/home/edit"
            className="inline-flex items-center gap-1.5 text-sm text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100 border border-neutral-300 dark:border-neutral-700 px-3 py-1.5 rounded-md hover:bg-neutral-100 dark:hover:bg-neutral-800 transition-colors"
          >
            <Pencil size={14} />
            编辑
          </Link>
        </div>
      )}
      <section className={`space-y-6 ${user ? 'pt-12' : 'pt-4'}`}>
        <div className="prose prose-neutral dark:prose-invert min-w-full">
          <ReactMarkdown
            remarkPlugins={[remarkGfm, remarkBreaks]}
            components={{
              table: ({ node, ...props }) => (
                <div className="w-full overflow-x-auto my-6 border border-neutral-200 dark:border-neutral-800 rounded-lg">
                  <table {...props} className="min-w-full divide-y divide-neutral-200 dark:divide-neutral-800" />
                </div>
              ),
            }}
          >
            {homeContent}
          </ReactMarkdown>
        </div>
      </section>

      <section className="pt-6 border-t border-neutral-200 dark:border-neutral-800">
        <div className="flex gap-6 items-center">
          <span className="text-sm font-medium text-neutral-500 dark:text-neutral-400">Get in touch</span>
          <div className="flex gap-4 items-center">
            {siteConfig.contact_github && (
              <a
                href={siteConfig.contact_github}
                target="_blank"
                rel="noopener noreferrer"
                className="text-neutral-500 hover:text-neutral-900 dark:text-neutral-400 dark:hover:text-neutral-100 transition-colors"
                title="GitHub"
                aria-label="GitHub"
              >
                <svg viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" strokeWidth="2" fill="none" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path>
                </svg>
              </a>
            )}
            {siteConfig.contact_twitter && (
              <a
                href={siteConfig.contact_twitter}
                target="_blank"
                rel="noopener noreferrer"
                className="text-neutral-500 hover:text-sky-500 dark:text-neutral-400 dark:hover:text-sky-400 transition-colors"
                title="Twitter"
                aria-label="Twitter"
              >
                <svg viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" strokeWidth="2" fill="none" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"></path>
                </svg>
              </a>
            )}
            {siteConfig.contact_mail && (
              <a
                href={`mailto:${siteConfig.contact_mail}`}
                className="text-neutral-500 hover:text-red-500 dark:text-neutral-400 dark:hover:text-red-400 transition-colors"
                title="Email"
                aria-label="Email"
              >
                <Mail size={20} />
              </a>
            )}
            {!siteConfig.contact_github && !siteConfig.contact_twitter && !siteConfig.contact_mail && (
              <span className="text-sm text-neutral-400 font-normal">（暂无联系方式）</span>
            )}
          </div>
        </div>
      </section>
    </div>
  )
}
