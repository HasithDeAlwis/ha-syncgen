"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowDown, Database, FileText, Zap, Shield, Code, Play, Activity, BarChart3 } from "lucide-react"
import Link from "next/link"

export default function HomePage() {
  const scrollToDemo = () => {
    document.getElementById("demo-section")?.scrollIntoView({ behavior: "smooth" })
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gradient-to-b from-background to-card py-20 px-4">
        <div className="container mx-auto max-w-4xl text-center">
          <div className="mb-8">
            <h1 className="text-4xl md:text-6xl font-bold text-foreground mb-6 text-balance">
              Scaffold High-Availability Databases with <span className="text-primary">YAML</span>
            </h1>
            <p className="text-xl text-muted-foreground mb-8 text-pretty max-w-2xl mx-auto">
              Automate the creation and setup of HA databases using declarative YAML configuration. Built for
              reliability, designed for DevOps.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-12">
            <Button asChild size="lg" className="text-lg px-8">
              <Link href="/docs">
                <FileText className="mr-2 h-5 w-5" />
                See Docs
              </Link>
            </Button>
            <Button
              variant="outline"
              size="lg"
              className="text-lg px-8 bg-transparent hover:bg-primary/10 transition-all duration-300"
              onClick={scrollToDemo}
            >
              See Demo
              <ArrowDown className="ml-2 h-5 w-5 animate-bounce" />
            </Button>
          </div>

          <div className="text-sm text-muted-foreground">Open source • GPL-3.0 License • Built with Go</div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 px-4">
        <div className="container mx-auto max-w-6xl">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Why ha-syncgen?</h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Simplify and standardize your database infrastructure with configuration-driven automation.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            <Card className="border-border">
              <CardHeader>
                <Activity className="h-12 w-12 text-primary mb-4" />
                <CardTitle>Built-in Health Checks</CardTitle>
                <CardDescription>
                  Continuous monitoring and automated health checks ensure your database infrastructure stays healthy
                  and responsive.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-border">
              <CardHeader>
                <Code className="h-12 w-12 text-primary mb-4" />
                <CardTitle>YAML Configuration</CardTitle>
                <CardDescription>
                  Define your database schemas, replication strategies, and HA configurations declaratively with simple
                  YAML files.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-border">
              <CardHeader>
                <Zap className="h-12 w-12 text-primary mb-4" />
                <CardTitle>Automated Scaffolding</CardTitle>
                <CardDescription>
                  Generate and deploy complete database infrastructures automatically, reducing manual setup time and
                  errors.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-border">
              <CardHeader>
                <Shield className="h-12 w-12 text-primary mb-4" />
                <CardTitle>Production Ready</CardTitle>
                <CardDescription>
                  Uses built-in PostgreSQL features like WAL streaming and pgbasebackup for reliable replication and
                  backup strategies.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-border">
              <CardHeader>
                <BarChart3 className="h-12 w-12 text-primary mb-4" />
                <CardTitle>Built-in Monitoring</CardTitle>
                <CardDescription>
                  Integrated monitoring with Datadog support provides real-time insights into your database performance
                  and health.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-border">
              <CardHeader>
                <Database className="h-12 w-12 text-primary mb-4" />
                <CardTitle>Open Source</CardTitle>
                <CardDescription>
                  GPL-3.0 licensed and open for community contributions. Extend and adapt for your specific needs.
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </div>
      </section>

      <section id="demo-section" className="py-20 px-4 bg-card/50">
        <div className="container mx-auto max-w-4xl">
          <div className="text-center mb-12">
            <div className="inline-flex items-center gap-2 bg-primary/10 text-primary px-4 py-2 rounded-full text-sm font-medium mb-4">
              <Play className="h-4 w-4" />
              Live Demo
            </div>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">See ha-syncgen in Action</h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Watch how easy it is to scaffold a high-availability database infrastructure with just a YAML file.
            </p>
          </div>

          <div className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-primary/20 to-secondary/20 rounded-lg blur opacity-75 group-hover:opacity-100 transition duration-1000 group-hover:duration-200 animate-tilt"></div>
            <div className="relative">
              <Card className="border-border overflow-hidden">
                <CardContent className="p-0">
                  <div className="aspect-video bg-muted/50 flex items-center justify-center">
                    <iframe
                      className="w-full h-full"
                      src="https://www.youtube.com/embed/dQw4w9WgXcQ"
                      title="ha-syncgen Demo"
                      frameBorder="0"
                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                      allowFullScreen
                    ></iframe>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>

          <div className="text-center mt-8">
            <p className="text-sm text-muted-foreground mb-4">
              Want to try it yourself? Get started with the documentation.
            </p>
            <Button asChild>
              <Link href="/docs">
                Get Started Now
                <ArrowDown className="ml-2 h-4 w-4 rotate-[-90deg]" />
              </Link>
            </Button>
          </div>
        </div>
      </section>

      {/* YAML Example Section */}
      <section className="py-20 px-4 bg-card">
        <div className="container mx-auto max-w-4xl">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Simple Configuration</h2>
            <p className="text-lg text-muted-foreground">
              Define your high-availability database setup with clean, readable YAML.
            </p>
          </div>

          <Card className="bg-background border-border">
            <CardContent className="p-6">
              <pre className="text-sm text-foreground overflow-x-auto">
                <code>{`# ha-syncgen.yaml
primary:
    host: 18.189.194.40
    port: 5432
    data_directory: /var/lib/pgsql/data
    db_name: primary_db
    db_user: admin
    db_password: admin_password
    replication_user: primary_replica_user
    replication_password: replica_user_password
replicas:
    - host: 18.118.139.199
      port: 5432
      replication_slot: replica_slot_1
      sync_mode: async
      db_user: replica1_admin
      db_password: replica1_password
options:
    promote_on_failure: true
    wal_level: replica
    max_wal_senders: 3
    wal_keep_size: 1GB
    hot_standby: true
    synchronous_commit: "on"
monitoring:
    datadog:
        enabled: true
        api_key: datadog_api_key
        site: datadoghq.com
        datadog_user_password: admin_password

`}</code>
              </pre>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4">
        <div className="container mx-auto max-w-2xl text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-6">Ready to Get Started?</h2>
          <p className="text-lg text-muted-foreground mb-8">
            Join developers and DevOps teams using ha-syncgen to build reliable database infrastructure.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button asChild size="lg" className="text-lg px-8">
              <Link href="/docs">Get Started</Link>
            </Button>
            <Button variant="outline" size="lg" className="text-lg px-8 bg-transparent" asChild>
              <Link href="https://github.com/HasithDeAlwis/ha-syncgen" target="_blank" rel="noopener noreferrer">
                View on GitHub
              </Link>
            </Button>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border py-12 px-4">
        <div className="container mx-auto max-w-4xl">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="mb-4 md:mb-0">
              <h3 className="text-lg font-semibold text-foreground">ha-syncgen</h3>
              <p className="text-sm text-muted-foreground">High-Availability Database Scaffolding</p>
            </div>
            <div className="flex space-x-6 text-sm text-muted-foreground">
              <Link href="/docs" className="hover:text-foreground transition-colors">
                Documentation
              </Link>
              <Link
                href="https://github.com/HasithDeAlwis/ha-syncgen"
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-foreground transition-colors"
              >
                GitHub
              </Link>
              <Link
                href="https://github.com/HasithDeAlwis/ha-syncgen/issues"
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-foreground transition-colors"
              >
                Issues
              </Link>
            </div>
          </div>
          <div className="mt-8 pt-8 border-t border-border text-center text-sm text-muted-foreground">
            © 2024 ha-syncgen. Licensed under GPL-3.0.
          </div>
        </div>
      </footer>
    </div>
  )
}
