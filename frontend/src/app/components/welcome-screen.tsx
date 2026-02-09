import * as React from 'react';
import { motion } from 'motion/react';
import { LaurelWreath, GreekPattern } from './greek-pattern';

interface WelcomeScreenProps {
  onEnter: () => void;
}

export function WelcomeScreen({ onEnter }: WelcomeScreenProps) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background flex items-center justify-center p-4 relative overflow-hidden">
      {/* Background decorative elements */}
      <div className="absolute inset-0 opacity-30">
        <GreekPattern className="absolute top-0 left-0 w-full h-8 text-primary" />
        <GreekPattern className="absolute bottom-0 left-0 w-full h-8 text-primary rotate-180" />
      </div>

      {/* Marble texture overlay */}
      <div className="absolute inset-0 opacity-5" 
           style={{
             backgroundImage: 'url("data:image/svg+xml,%3Csvg width=\'100\' height=\'100\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cfilter id=\'noise\'%3E%3CfeTurbulence baseFrequency=\'0.9\' numOctaves=\'4\'/%3E%3C/filter%3E%3Crect width=\'100\' height=\'100\' filter=\'url(%23noise)\' opacity=\'0.5\'/%3E%3C/svg%3E")',
           }} 
      />

      <motion.div
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.8, ease: 'easeOut' }}
        className="relative z-10 max-w-2xl w-full"
      >
        <div className="bg-card/80 backdrop-blur-sm border-2 border-primary/30 rounded-xl shadow-2xl p-8 md:p-12 relative">
          {/* Decorative laurel */}
          <div className="absolute -top-16 left-1/2 -translate-x-1/2">
            <LaurelWreath className="w-32 h-32 text-primary" />
          </div>

          <div className="text-center space-y-6">
            {/* Logo and Title */}
            <motion.div
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3, duration: 0.6 }}
            >
              <h1 className="text-5xl md:text-6xl font-bold text-primary tracking-wider mb-2"
                  style={{ fontFamily: 'var(--font-heading)' }}>
                ΠΥΘΙΑ
              </h1>
              <h2 className="text-2xl md:text-3xl text-foreground/80"
                  style={{ fontFamily: 'var(--font-heading)' }}>
                Pythia
              </h2>
              <p className="text-base md:text-lg text-muted-foreground mt-2 italic">
                The Oracle of The Language
              </p>
            </motion.div>

            {/* Description */}
            <motion.p
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.6, duration: 0.6 }}
              className="text-foreground/70 max-w-md mx-auto leading-relaxed"
            >
              Embark on your sacred journey of language mastery. 
              Let the Oracle guide you through focused learning sessions, 
              one word at a time.
            </motion.p>

            {/* Enter Button */}
            <motion.button
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.9, duration: 0.6 }}
              onClick={onEnter}
              className="mt-8 px-12 py-4 bg-gradient-to-r from-primary to-accent text-primary-foreground 
                         rounded-lg font-semibold text-lg shadow-lg hover:shadow-xl 
                         transform hover:scale-105 transition-all duration-300
                         border border-primary/20 relative overflow-hidden group"
            >
              <span className="relative z-10">Enter the Oracle</span>
              <div className="absolute inset-0 bg-gradient-to-r from-accent to-primary opacity-0 
                              group-hover:opacity-100 transition-opacity duration-300" />
            </motion.button>

            {/* Decorative quote */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 1.2, duration: 0.6 }}
              className="mt-8 pt-6 border-t border-primary/20"
            >
              <p className="text-sm text-muted-foreground italic">
                "γνῶθι σεαυτόν"
              </p>
              <p className="text-xs text-muted-foreground mt-1">
                Know Thyself
              </p>
            </motion.div>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
