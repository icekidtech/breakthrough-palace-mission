'use client'

import Link from 'next/link'
import { useState } from 'react'
import {
  ArrowUpRight,
  Bell,
  BookOpen,
  CalendarDays,
  Check,
  ChevronDown,
  CircleDollarSign,
  FileText,
  HeartHandshake,
  LayoutDashboard,
  Menu,
  MoreHorizontal,
  Search,
  Settings,
  ShieldCheck,
  Users,
  X,
} from 'lucide-react'

const navItems = [
  { label: 'Overview', icon: LayoutDashboard },
  { label: 'Family roster', icon: Users },
  { label: 'Meetings & minutes', icon: FileText },
  { label: 'Attendance', icon: CalendarDays },
  { label: 'Contributions', icon: CircleDollarSign },
]

const meetings = [
  { date: '24', month: 'OCT', title: 'Family thanksgiving meeting', meta: 'Sat, 24 Oct · 4:00 PM', status: 'Upcoming', tone: 'blue' },
  { date: '10', month: 'OCT', title: 'Monthly family meeting', meta: 'Sat, 10 Oct · 4:00 PM', status: 'Minutes published', tone: 'green' },
  { date: '26', month: 'SEP', title: 'September family meeting', meta: 'Sat, 26 Sep · 4:00 PM', status: 'Minutes published', tone: 'slate' },
]

const activities = [
  { icon: Check, title: 'Attendance recorded', description: 'September family meeting', time: '2 hours ago', color: 'green' },
  { icon: FileText, title: 'Minutes published', description: 'Monthly family meeting', time: 'Yesterday', color: 'blue' },
  { icon: CircleDollarSign, title: 'Contribution confirmed', description: 'Dominion family dues', time: '2 days ago', color: 'amber' },
]

export function FamilyDashboard() {
  const [active, setActive] = useState('Overview')
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <div className="min-h-screen bg-[#f4f8fc] text-[#172a46]">
      <aside className={`fixed inset-y-0 left-0 z-30 flex w-[256px] flex-col bg-[#132745] px-5 py-6 text-white transition-transform lg:translate-x-0 ${mobileOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div className="flex items-center justify-between px-2">
          <div className="flex items-center gap-3">
            <div className="flex size-10 items-center justify-center rounded-[13px] bg-[#8dd6ff] text-[#132745] shadow-[0_8px_24px_rgba(141,214,255,0.24)]">
              <HeartHandshake className="size-5" strokeWidth={2.4} />
            </div>
            <div>
              <p className="text-[15px] font-bold tracking-wide">Breakthrough</p>
              <p className="text-[10px] font-medium uppercase tracking-[0.18em] text-[#8da8c8]">Palace Mission</p>
            </div>
          </div>
          <button aria-label="Close navigation" className="rounded-lg p-2 text-[#a9bfd8] lg:hidden" onClick={() => setMobileOpen(false)}><X /></button>
        </div>

        <div className="mt-12 flex flex-1 flex-col">
          <p className="px-3 text-[10px] font-semibold uppercase tracking-[0.2em] text-[#7592b4]">Workspace</p>
          <nav className="mt-4 flex flex-col gap-1.5" aria-label="Main navigation">
            {navItems.map((item) => {
              const Icon = item.icon
              const selected = active === item.label
              return item.label === 'Meetings & minutes' ? <Link key={item.label} href="/meetings" onClick={() => setMobileOpen(false)} className={`flex items-center gap-3 rounded-xl px-3 py-3 text-left text-[13px] font-medium transition-colors ${selected ? 'bg-[#2e75c7] text-white shadow-[0_8px_18px_rgba(46,117,199,0.28)]' : 'text-[#a9bfd8] hover:bg-white/8 hover:text-white'}`}><Icon className="size-[18px]" />{item.label}</Link> : <button key={item.label} onClick={() => { setActive(item.label); setMobileOpen(false) }} className={`flex items-center gap-3 rounded-xl px-3 py-3 text-left text-[13px] font-medium transition-colors ${selected ? 'bg-[#2e75c7] text-white shadow-[0_8px_18px_rgba(46,117,199,0.28)]' : 'text-[#a9bfd8] hover:bg-white/8 hover:text-white'}`}><Icon className="size-[18px]" />{item.label}</button>
            })}
          </nav>
          <p className="mt-10 px-3 text-[10px] font-semibold uppercase tracking-[0.2em] text-[#7592b4]">Account</p>
          <nav className="mt-4 flex flex-col gap-1.5">
            <Link href="/settings" className="flex items-center gap-3 rounded-xl px-3 py-3 text-left text-[13px] font-medium text-[#a9bfd8] hover:bg-white/8 hover:text-white"><Settings className="size-[18px]" />Settings</Link>
            <Link href="/privacy-access" className="flex items-center gap-3 rounded-xl px-3 py-3 text-left text-[13px] font-medium text-[#a9bfd8] hover:bg-white/8 hover:text-white"><ShieldCheck className="size-[18px]" />Privacy & access</Link>
          </nav>
        </div>
        <div className="rounded-2xl border border-white/10 bg-white/6 p-4">
          <div className="flex items-center gap-3"><div className="flex size-9 items-center justify-center rounded-full bg-[#d9efff] text-sm font-bold text-[#2764a3]">DA</div><div className="min-w-0"><p className="truncate text-xs font-semibold">Daniel Adeyemi</p><p className="mt-0.5 text-[11px] text-[#8da8c8]">Member · Dominion</p></div><MoreHorizontal className="ml-auto size-4 text-[#8da8c8]" /></div>
        </div>
      </aside>

      {mobileOpen && <button aria-label="Close navigation overlay" className="fixed inset-0 z-20 bg-[#0e2039]/40 backdrop-blur-sm lg:hidden" onClick={() => setMobileOpen(false)} />}

      <main className="lg:pl-[256px]">
        <header className="flex h-[76px] items-center justify-between border-b border-[#dce8f3] bg-white/80 px-5 backdrop-blur-xl sm:px-8 lg:px-10">
          <div className="flex items-center gap-3"><button aria-label="Open navigation" className="rounded-xl p-2 text-[#426487] hover:bg-[#edf5fb] lg:hidden" onClick={() => setMobileOpen(true)}><Menu /></button><div className="hidden items-center gap-2 text-xs text-[#7690aa] sm:flex"><span>Family workspace</span><span>/</span><span className="font-medium text-[#274463]">Overview</span></div><span className="text-sm font-semibold text-[#1f3857] sm:hidden">Overview</span></div>
          <div className="flex items-center gap-2 sm:gap-4"><button aria-label="Search" className="rounded-xl p-2.5 text-[#55718e] hover:bg-[#edf5fb]"><Search className="size-[19px]" /></button><button aria-label="Notifications" className="relative rounded-xl p-2.5 text-[#55718e] hover:bg-[#edf5fb]"><Bell className="size-[19px]" /><span className="absolute right-2 top-2 size-1.5 rounded-full bg-[#2f8ad8]" /></button><div className="hidden h-7 w-px bg-[#e3edf5] sm:block" /><button className="flex items-center gap-2 rounded-xl p-1.5 pr-2 hover:bg-[#edf5fb]"><div className="flex size-8 items-center justify-center rounded-full bg-[#d9efff] text-xs font-bold text-[#2764a3]">DA</div><ChevronDown className="size-4 text-[#7792ad]" /></button></div>
        </header>

        <div className="mx-auto max-w-[1380px] px-5 py-7 sm:px-8 lg:px-10 lg:py-9">
          <section className="flex flex-col justify-between gap-5 md:flex-row md:items-end"><div><p className="text-sm font-medium text-[#5d7b99]">Saturday, 4 October 2026</p><h1 className="mt-2 text-3xl font-bold tracking-[-0.04em] text-[#142e4d] sm:text-[38px]">Good morning, Daniel <span className="text-[#4295d5]">.</span></h1><p className="mt-2 max-w-xl text-sm leading-6 text-[#7290ad]">Here&apos;s what&apos;s happening in <span className="font-semibold text-[#41627f]">Dominion Family</span> this week.</p></div><button className="flex w-fit items-center gap-2 rounded-xl bg-[#2e75c7] px-4 py-3 text-sm font-semibold text-white shadow-[0_10px_20px_rgba(46,117,199,0.2)] hover:bg-[#2569b7]"><CalendarDays className="size-4" />View calendar</button></section>

          <section className="mt-8 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            <StatCard label="Family members" value="48" detail="2 new this month" icon={Users} tone="blue" />
            <StatCard label="Attendance rate" value="86%" detail="+4.2% from last month" icon={CalendarDays} tone="green" />
            <StatCard label="My contributions" value="₦85,000" detail="₦15,000 pending" icon={CircleDollarSign} tone="amber" />
            <StatCard label="Published minutes" value="12" detail="All caught up" icon={BookOpen} tone="purple" />
          </section>

          <section className="mt-6 grid gap-6 xl:grid-cols-[minmax(0,1.65fr)_minmax(320px,0.9fr)]">
            <div className="rounded-[22px] border border-[#dce8f3] bg-white p-5 shadow-[0_12px_36px_rgba(49,92,130,0.05)] sm:p-6"><div className="flex items-center justify-between"><div><h2 className="text-base font-bold text-[#1a3858]">Upcoming meetings</h2><p className="mt-1 text-xs text-[#89a0b6]">Stay connected with your family</p></div><button className="text-xs font-semibold text-[#2d79c8] hover:underline">See all meetings</button></div><div className="mt-5 flex flex-col gap-3">{meetings.map((meeting) => <MeetingRow key={meeting.title} {...meeting} />)}</div></div>
            <div className="rounded-[22px] bg-[#2776c8] p-6 text-white shadow-[0_16px_35px_rgba(39,118,200,0.2)]"><div className="flex items-start justify-between"><div><p className="text-xs font-medium text-[#b9ddfa]">Your family contribution</p><h2 className="mt-3 text-[30px] font-bold tracking-[-0.04em]">₦85,000</h2><p className="mt-1 text-xs text-[#c4e2f8]">Total confirmed this year</p></div><div className="flex size-10 items-center justify-center rounded-xl bg-white/15"><CircleDollarSign className="size-5" /></div></div><div className="mt-7"><div className="mb-2 flex justify-between text-[11px] text-[#c6e4fa]"><span>Annual family dues</span><span>68%</span></div><div className="h-2 rounded-full bg-white/20"><div className="h-2 w-[68%] rounded-full bg-white" /></div></div><button className="mt-7 flex items-center gap-2 text-xs font-semibold text-white hover:text-[#d8efff]">View contribution history <ArrowUpRight className="size-4" /></button></div>
          </section>

          <section className="mt-6 grid gap-6 xl:grid-cols-[minmax(0,1fr)_minmax(320px,0.72fr)]">
            <div className="rounded-[22px] border border-[#dce8f3] bg-white p-5 shadow-[0_12px_36px_rgba(49,92,130,0.05)] sm:p-6"><div className="flex items-center justify-between"><div><h2 className="text-base font-bold text-[#1a3858]">Attendance overview</h2><p className="mt-1 text-xs text-[#89a0b6]">Your participation over the last 6 months</p></div><button className="flex items-center gap-1 rounded-lg border border-[#dce8f3] px-2.5 py-2 text-xs font-medium text-[#58738f]">Last 6 months <ChevronDown className="size-3.5" /></button></div><div className="mt-6 flex h-[168px] items-end gap-3 border-b border-[#e9f0f6] px-2 sm:gap-6">{[68, 75, 64, 88, 79, 94].map((height, i) => <div key={i} className="flex h-full flex-1 flex-col items-center justify-end gap-2"><div className={`w-full max-w-[38px] rounded-t-lg ${i === 5 ? 'bg-[#2e75c7]' : 'bg-[#cce7f8]'}`} style={{ height: `${height}%` }} /><span className="text-[10px] font-medium text-[#93a8ba]">{['May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'][i]}</span></div>)}</div><div className="mt-4 flex items-center gap-5 text-[11px] text-[#7892ab]"><span className="flex items-center gap-2"><i className="size-2 rounded-full bg-[#2e75c7]" />Present</span><span className="flex items-center gap-2"><i className="size-2 rounded-full bg-[#cce7f8]" />Other</span><span className="ml-auto font-semibold text-[#2e75c7]">86% attendance rate</span></div></div>
            <div className="rounded-[22px] border border-[#dce8f3] bg-white p-5 shadow-[0_12px_36px_rgba(49,92,130,0.05)] sm:p-6"><div className="flex items-center justify-between"><div><h2 className="text-base font-bold text-[#1a3858]">Recent activity</h2><p className="mt-1 text-xs text-[#89a0b6]">Latest updates from your family</p></div><button aria-label="More activity options" className="rounded-lg p-1.5 text-[#7892ab] hover:bg-[#eef5fa]"><MoreHorizontal className="size-5" /></button></div><div className="mt-5 flex flex-col gap-5">{activities.map((activity) => <ActivityRow key={activity.title} {...activity} />)}</div></div>
          </section>

          <section className="mt-6 flex flex-col gap-4 rounded-[22px] border border-[#dce8f3] bg-[#eef8ff] p-5 sm:flex-row sm:items-center sm:justify-between sm:p-6"><div className="flex items-start gap-4"><div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-white text-[#2e75c7] shadow-sm"><Bell className="size-5" /></div><div><h2 className="text-sm font-bold text-[#1d456b]">Family announcement</h2><p className="mt-1 text-xs leading-5 text-[#64839f]">Our family thanksgiving service is coming up on 24 October. Let&apos;s come together in gratitude.</p></div></div><button className="whitespace-nowrap text-xs font-semibold text-[#2d79c8] hover:underline">Read announcement <ArrowUpRight className="ml-1 inline size-3.5" /></button></section>
          <p className="mt-8 pb-3 text-center text-[11px] text-[#9ab0c3]">Dominion Family · Faithful records platform</p>
        </div>
      </main>
    </div>
  )
}

function StatCard({ label, value, detail, icon: Icon, tone }: { label: string; value: string; detail: string; icon: typeof Users; tone: string }) {
  const tones: Record<string, string> = { blue: 'bg-[#e8f5ff] text-[#2e75c7]', green: 'bg-[#e8f8f0] text-[#249563]', amber: 'bg-[#fff6e4] text-[#bf8226]', purple: 'bg-[#f0edff] text-[#735bca]' }
  return <div className="rounded-[18px] border border-[#dce8f3] bg-white p-5 shadow-[0_10px_30px_rgba(49,92,130,0.04)]"><div className="flex items-center justify-between"><p className="text-xs font-medium text-[#7892ab]">{label}</p><div className={`flex size-9 items-center justify-center rounded-xl ${tones[tone]}`}><Icon className="size-[18px]" /></div></div><p className="mt-4 text-[25px] font-bold tracking-[-0.04em] text-[#183654]">{value}</p><p className="mt-1 text-[11px] font-medium text-[#5e9a7d]">{detail}</p></div>
}

function MeetingRow({ date, month, title, meta, status, tone }: { date: string; month: string; title: string; meta: string; status: string; tone: string }) {
  return <div className="flex items-center gap-3 rounded-2xl border border-[#e7eff6] p-3.5 transition-colors hover:border-[#b9dbf2] hover:bg-[#f8fcff] sm:gap-4"><div className="flex size-11 shrink-0 flex-col items-center justify-center rounded-xl bg-[#edf7ff] text-[#2c76c4]"><span className="text-base font-bold leading-none">{date}</span><span className="mt-1 text-[9px] font-bold tracking-wider">{month}</span></div><div className="min-w-0 flex-1"><p className="truncate text-sm font-semibold text-[#274663]">{title}</p><p className="mt-1 text-[11px] text-[#8aa0b5]">{meta}</p></div><span className={`hidden rounded-full px-2.5 py-1.5 text-[10px] font-semibold sm:block ${tone === 'blue' ? 'bg-[#e9f5ff] text-[#2d78c6]' : tone === 'green' ? 'bg-[#e9f8f0] text-[#2c9565]' : 'bg-[#f1f4f7] text-[#7790a7]'}`}>{status}</span><MoreHorizontal className="size-4 shrink-0 text-[#a0b2c2]" /></div>
}

function ActivityRow({ icon: Icon, title, description, time, color }: { icon: typeof Check; title: string; description: string; time: string; color: string }) {
  const colors: Record<string, string> = { green: 'bg-[#e8f8f0] text-[#249563]', blue: 'bg-[#e8f5ff] text-[#2e75c7]', amber: 'bg-[#fff6e4] text-[#bf8226]' }
  return <div className="flex items-center gap-3"><div className={`flex size-9 shrink-0 items-center justify-center rounded-xl ${colors[color]}`}><Icon className="size-4" /></div><div className="min-w-0 flex-1"><p className="truncate text-xs font-semibold text-[#385773]">{title}</p><p className="mt-0.5 truncate text-[11px] text-[#8ba1b5]">{description}</p></div><span className="shrink-0 text-[10px] text-[#9aaebe]">{time}</span></div>
}

export default FamilyDashboard

